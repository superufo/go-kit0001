package main

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type Data []byte

// wsn writing Serial Number   rsn   Reading Serial Number
type DataFile interface {
	//读取一个数据块文件
	Read() (rsn int64, d Data, err error)
	//写入一个数据块
	Write(d Data) (wsn uint64, err error)
	//获取最后读取的的最大序列号
	RSn()
	//获取最后写入的最大序列好
	WSN()
	//获取数据块的长度
	dataLen() uint32
	//关闭文件
	Close() error
}

/**
当有多个写操作同时要增加woffset字段的时候，会产生竟态条件，需要互斥锁Wmutex来加以保护
同理 互斥锁Rmutex用来消除多个读增加　ｒｏｆｆｓｅｔ
 */
type myDataFile struct {
	f       *os.File     //文件句柄
	fmutex  sync.RWMutex //用于文件的读写锁
	woffset int64        //写操作需要用到的偏移量
	roffset uint64       //读操作用到的偏移量
	wmutex  sync.Mutex   //写操作用到的互斥量
	rmutex  sync.Mutex   //读操作用到的互斥量
	dataLen uint32       //数据块长度
	rcond   *sync.Cond    //条件变量
}

func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if dataLen == 0 {
		return nil, errors.New("Invide Data Length")
	}

	// 未初始化的字段是0 值
	df := &myDataFile{f: f, dataLen: dataLen}
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

/***
* rmutex 保证了 offset = int64(df.roffset)  df.roffset += uint64(df.dataLen)
* 代码的执行是互斥的, 多个读操作快于写操作，使得ReadAt无数据可读，返回io.EOF
* io.EOF (代表无数据可读)
* 问题 调用方会读取到出错的数据块的序列号,但无法再次读取到尝试读取这个数据块
* 由于其他正在或后续执行的Read方法会继续增加读偏移量roffset的值，因此当调用
× 再次调用这个Read方法的时候，只可能读取到在此数据块后面的数据块
*/
func (df *myDataFile) Read1() (rsn int64, err error) {
	//读取文件更新偏移量
	var offset int64
	df.rmutex.Lock()
	offset = int64(df.roffset)
	df.roffset += uint64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块
	rsn = offset / int64(df.dataLen)
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	bytes := make([]byte, df.dataLen)
	_, err = df.f.ReadAt(bytes, offset)
	if err != nil {
		return
	}

	return
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	//读取文件更新偏移量
	var offset int64
	df.rmutex.Lock()
	offset = int64(df.roffset)
	df.roffset += uint64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	for {
		df.fmutex.RLock()
		//defer df.fmutex.RUnlock()
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				//当出现io.EOF 让当前的Goroutine放弃当前的fmutex的读锁，并且等待通知的到来，
				//放弃当前的读锁，意味Write 操作的写操作不会受到它的阻碍
				//一旦有新的写操作完成，应该及时向条件变量发送通知，已唤醒等待的Goroutine
				//一旦唤醒等待的Goroutine，再次检查满足的条件
				df.rcond.Wait()
				continue
			}

			df.fmutex.RUnlock()
			return
		}

		d = bytes
		df.fmutex.RUnlock()
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	//读取并更新写偏移量
	var offset int64
	for {
		//原子操作支持的6种数据 int32 int64  unit32 uint64  uintptr 和 unsaft.Pointer
		//原子操作 增和减 比较并且交换 载入 存储 和 交换
		offset = atomic.LoadInt64(&df.woffset)
		if atomic.CompareAndSwapInt64(&df.woffset, offset, (offset + int64(df.dataLen))) {
			break
		}
	}

	//写入一个数据块
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	//当参数d 的上都大于数据块的最大长度 会先进行截短处理再将数据写入文件　
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	//写入文件
	_, err = df.f.Write(bytes)
	//唤醒条件变量
	df.rcond.Signal()

	return
}

func (df *myDataFile) RSN() int64 {
   offset := atomic.LoadInt64(&df.roffset)
   return offset / int64(df.dataLen)
}

func (df *myDataFile) WSN() int64 {
	offset := atomic.LoadInt64(&df.woffset)
	return offset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	if df.f == nil {
		return nil
	}
	return df.f.Close()
}