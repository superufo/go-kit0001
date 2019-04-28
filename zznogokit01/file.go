package main

import (
	"github.com/kataras/iris/core/errors"
	"io"
	"os"
	"sync"
)

type Data []byte

// wsn writing Serial Number   rsn   Reading Serial Number
type DataFile interface{
	//读取一个数据块文件
	Read() (rsn int64,d Data,err error)
	//写入一个数据块
	Write(d Data) (wsn uint64,err error)
	//获取最后读取的的最大序列号
	RSn()
	//获取最后写入的最大序列好
	WSN()
	//获取数据块的长度
	DataLen() uint32
	//关闭文件
	Close() error
}

/**
当有多个写操作同时要增加woffset字段的时候，会产生竟态条件，需要互斥锁Wmutex来加以保护
同理 互斥锁Rmutex用来消除多个读增加　ｒｏｆｆｓｅｔ
 */
type  myDataFile struct {
	f   *os.File               //文件句柄
	fmutex sync.RWMutex        //用于文件的读写锁
	woffset int64              //写操作需要用到的偏移量
	roffset uint64            //读操作用到的偏移量
	wmutex  sync.Mutex        //写操作用到的互斥量
	rmutex sync.Mutex         //读操作用到的互斥量
	datalen uint32            //数据块长度
}

func NewDataFile(path string,datalen uint32)(DataFile,error){
	f,err := os.Create(path)
	if err!=nil{
		return nil,err
	}

	if datalen==0{
		return nil,errors.New("Invide Data Length")
	}

	// 未初始化的字段是0 值
	df := &myDataFile{f:f,datalen:datalen}
	return df,nil
}

/***
* rmutex 保证了 offset = int64(df.roffset)  df.roffset += uint64(df.datalen)
* 代码的执行是互斥的, 多个读操作快于写操作，使得ReadAt无数据可读，返回io.EOF
* io.EOF (代表无数据可读)
* 问题 调用方会读取到出错的数据块的序列号,但无法再次读取到尝试读取这个数据块
* 由于其他正在或后续执行的Read方法会继续增加读偏移量roffset的值，因此当调用
× 再次调用这个Read方法的时候，只可能读取到在此数据块后面的数据块
*/
func (df *myDataFile) Read1()(rsn int64,err error){
	//读取文件更新偏移量
	var offset int64
	df.rmutex.Lock()
	offset = int64(df.roffset)
	df.roffset += uint64(df.datalen)
	df.rmutex.Unlock()

	//读取一个数据块
	rsn =offset/int64(df.datalen)
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	bytes := make([]byte,df.datalen)
	_,err =df.f.ReadAt(bytes,offset)
	if err != nil {
		return
	}

	return
}

func (df *myDataFile) Read()(rsn int64,d Data,err error){
	//读取文件更新偏移量
	var offset int64
	df.rmutex.Lock()
	offset = int64(df.roffset)
	df.roffset += uint64(df.datalen)
	df.rmutex.Unlock()

	//读取一个数据块
	rsn =offset/int64(df.datalen)
	bytes := make([]byte,df.datalen)
	for {
		df.fmutex.RLock()
		//defer df.fmutex.RUnlock()
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err==io.EOF {
				df.fmutex.RUnlock()
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

func (df *myDataFile) Write(d Data) (wsn int64,err error){
	//读取并更新写偏移量
	var offset int64
	df.wmutex.Lock()
	offset = df.woffset
	df.woffset += int64(df.datalen)

	//写入一个数据块
	wsn =offset/int64(df.datalen)
	var bytes []byte
	//当参数d 的上都大于数据块的最大长度 会先进行截短处理再将数据写入文件　
	if len(d) > int(df.datalen){
		bytes = d[0:df.datalen]
	}else{
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	//写入文件
	_,err = df.f.Write(bytes)

	return
}

func (df *myDataFile)RSN() int64{
	df.rmutex.Lock()
	defer df.rmutex.Unlock()
	return int64(df.roffset)/int64(df.datalen)
}

func (df *myDataFile) WSN() int64{
	df.wmutex.Lock()
	defer df.wmutex.Unlock()
	return df.woffset/int64(df.datalen)
}

