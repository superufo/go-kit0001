package main


import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"google.golang.org/grpc"

	"discovery/book"
)

/*
* 1、 连接注册中心
* 2、 获取提供的服务
* 3、 监听服务目录变化，目录变化更新本地缓存
* 4、 创建负载均衡器
* 5、 获取请求的 endPoint
*/

func main() {
	fmt.Println("start........")
	var (
		etcdServer = "47.112.111.171:2379"

		// 监听的服务前缀
		prefix     = "/services/book"
		ctx        = context.Background()
	)

	options := etcdv3.ClientOptions{
		DialTimeout:time.Second*100,
		DialKeepAlive:time.Second*100,
	}

	// 连接注册中心
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil{
		panic(err)
	}

	logger := log.NewNopLogger()
	//key := prefix+"/etcd001"
	// 创建实例管理器, 此管理器会 Watch 监听 etc 中 prefix 的目录变化更新缓存的服务实例数据  prefix
	instance,err := etcdv3.NewInstancer(client,prefix,logger)
	if err!=nil {
		panic(err)
	}
	fmt.Printf("instance:%+v\n",instance)

	// 创建端点管理器， 此管理器根据 Factory 和监听的到实例创建 endPoint
	// 并订阅 instancer 的变化动态更新 Factory 创建的 endPoint
	endpointer := sd.NewEndpointer(instance,reqFactory,logger)

	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	/**
	  我们可以通过负载均衡器直接获取请求的 endPoint，发起请求
	  reqEndPoint,_ := balancer.Endpoint()
	  */

	/**
	也可以通过 retry 定义尝试次数进行请求
	*/

	reqEndPoint := lb.Retry(300,3*time.Second,balancer)
	fmt.Printf("reqEndPoint:%+v\n",reqEndPoint)
	req := struct {}{}

	if _,err = reqEndPoint(ctx,req); err != nil {
		panic(err)
	}
	fmt.Println("end........")
}

// 通过传入的 实例地址  创建对应的请求 endPoint
func reqFactory(instanceAddr string)(endpoint.Endpoint,io.Closer,error){
	return func(ctx context.Context,request interface{})(interface{},error){
		fmt.Println("请求服务",instanceAddr)
		conn,err := grpc.Dial(instanceAddr,grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer conn.Close()

		bookClient := book.NewBookServiceClient(conn)
		bi,_:= bookClient.GetBookInfo(context.Background(),&book.BookInfoParams{BookId:1})
		fmt.Println("获取书籍详情")
		fmt.Println("bookId: 1", "=>", "bookName:", bi.BookName)

		bl,_ := bookClient.GetBookList(context.Background(), &book.BookListParams{Page:1, Limit:10})
		fmt.Println("获取书籍列表")
		for _,b := range bl.BookList {
			fmt.Println("bookId:", b.BookId, "=>", "bookName:", b.BookName)
		}
		return nil,nil
	},nil,nil
}