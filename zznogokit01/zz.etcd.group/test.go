package zz_etcd_group

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

/**
etcdmain: etcd Version: 3.3.8
etcd1    | 2019-04-01 18:00:24.172320 I | etcdmain: Git SHA: 33245c6b5
etcd1    | 2019-04-01 18:00:24.172326 I | etcdmain: Go Version: go1.9.7
etcd1    | 2019-04-01 18:00:24.172329 I | etcdmain: Go OS/Arch: linux/amd64
 */
func main(){
	cli,err := clientv3.New(clientv3.Config{
		Endpoints : []string{"127.0.0.1:32771","127.0.0.1:32773","127.0.0.1:32769"},
        DialTimeout: time.Second*10,
	})

	//Grant：分配一个租约。
	//Revoke：释放一个租约。
	//TimeToLive：获取剩余TTL时间。
	//Leases：列举所有etcd中的租约。
	//KeepAlive：自动定时的续约某个租约。
	//KeepAliveOnce：为某个租约续约一次。
	//Close：貌似是关闭当前客户端建立的所有租约
	lease := clientv3.NewLease(cli)
	fmt.Printf("lease:%+v \n", lease)
	grantResp, err := lease.Grant(context.TODO(), 10)
	if err!=nil {
		fmt.Printf("putrep:%+v \n", grantResp)
	}

	kv := clientv3.NewKV(cli)

	res,err :=  kv.Put(context.TODO(), "/test/expireme", "gone...", clientv3.WithLease(grantResp.ID))
	if err!=nil {
		fmt.Printf("res:%+v \n  err %+v \n ", res,err)
	}

	//当我们实现服务注册时，需要主动给Lease进行续约，这需要调用KeepAlive/KeepAliveOnce，你可以在一个循环中定时的调用
	//KeepAlive和Put一样，如果在执行之前Lease就已经过期了，那么需要重新分配Lease
	 keepres,err := lease.KeepAlive(context.TODO(),grantResp.ID)
	if err!=nil {
		fmt.Printf("keepres:%+v \n  err %+v \n ", keepres,err)
	}


    putrep,err := kv.Put(context.TODO(),"/test/a","some goods")
    if err!=nil {
		fmt.Printf("putrep:%+v \n", putrep)
	}

	// 再写一个孩子
	kv.Put(context.TODO(),"/test/b", "another")

	// 再写一个同前缀的干扰项
	kv.Put(context.TODO(), "/testmmmmm", "干扰")

    getres,err := kv.Get(context.TODO(),"/test/a")
	if err!=nil {
		fmt.Printf("getres:%+v \n", getres)
	}

    //特别的Get选项，获取/test目录下的所有孩子 WithPrefix()是指查找以/test/为前缀的所有key，因此可以模拟出查找子目录的效果
    rangres,err :=kv.Get(context.TODO(),"/test",clientv3.WithPrefix())
	if err !=nil {
		fmt.Printf("rangres:%+v \n", rangres)
	}

    //Op是一个抽象的操作，可以是Put/Get/Delete…
	op1 := clientv3.OpPut("/hi", "hello", clientv3.WithPrevKV())
	opResp, err := kv.Do(context.TODO(), op1)
	if err!=nil {
		fmt.Printf("opResp:%+v \n  err %+v \n ", opResp,err)
	}

	/************事务**************/
	tran := kv.Txn(context.TODO())
	//Value(“/hi”)是指key=/hi对应的value
	tranResp, err := tran.If(clientv3.Compare(clientv3.Value("/hi"), "=", "hello")).
		Then(clientv3.OpGet("/hi")).
		Else(clientv3.OpGet("/test/", clientv3.WithPrefix())).
		Commit()
	if err!=nil {
		if tranResp.Succeeded { // If = true
			fmt.Println("~~~", tranResp.Responses[0].GetResponseRange().Kvs)
		} else { // If =false
			fmt.Println("!!!", tranResp.Responses[0].GetResponseRange().Kvs)
		}
	}


}