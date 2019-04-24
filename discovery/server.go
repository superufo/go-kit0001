package main

//cd gopro\src\github.com\go-kit\kit\zz.wt\
//cd D:\gopro\src\github.com\go-kit\kit\zz.wt\discovery
import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	grpc_transport "github.com/go-kit/kit/transport/grpc"

	"discovery/book"
)

type Balancer interface {
	Endpoint() (endpoint.Endpoint, error)
}

var (
	//etcd服务地址
	etcdServer = "47.112.111.171:2379"

	//服务的信息目录
	prefix = "/services/book"

	//当前IP
	ip = "127.0.0.1"

	//当前启动服务实例的地址
	instance = ip+":56555"

	//服务实例注册的路径
	key = prefix+"/"+instance

	//服务实例注册的Val
	value = instance
	ctx = context.Background()

	//服务监听地址
	serviceAddress = ":56555"
)

type BookServer struct{
	bookListHandler grpc_transport.Handler
	bookInfoHandler grpc_transport.Handler
}

func getExternal() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	return string(content)
}

//通过grpc调用GetBookInfo时,GetBookInfo只做数据透传,
//调用BookServer中对应Handler.ServeGRPC转交给go-kit处理
func (s *BookServer) GetBookInfo(ctx context.Context,in  *book.BookInfoParams) (*book.BookInfo,error) {
	_,rsp,err := s.bookInfoHandler.ServeGRPC(ctx,in)
	if err != nil {
		return nil,err
	}
	return rsp.(*book.BookInfo),err
}

//通过grpc调用GetBookList时,GetBookList只做数据透传,
//调用BookServer中对应Handler.ServeGRPC转交给go-kit处理
func (s *BookServer) GetBookList(ctx context.Context,in *book.BookListParams)(*book.BookList,error){
	_,rsp,err := s.bookListHandler.ServeGRPC(ctx,in)
	if err !=nil{
		return nil,err
	}

	return rsp.(*book.BookList),err
}

func makeGetBookListEndpoint() endpoint.Endpoint{
	return func(ctx context.Context,request interface{})(interface{},error){
		list := new(book.BookList)
		//正式系统 这里常常从database 或 redis 取数据
		list.BookList = append(list.BookList,&book.BookInfo{BookId:3,BookName:"室外装修大全"})
		list.BookList = append(list.BookList,&book.BookInfo{BookId:2,BookName:"射雕英雄传"})
		list.BookList = append(list.BookList,&book.BookInfo{BookId:1,BookName:"天龙八部"})
		fmt.Printf("book list:%+v\n",list)
		return list,nil
	}
}

func makeGetBookInfoEndpoint() endpoint.Endpoint{
	return func (ctx context.Context,request interface{})(interface{},error){
		fmt.Printf("*book.BookInfoParams:%+v\n",request.(*book.BookInfoParams))
		req := request.(*book.BookInfoParams)
		//正式系统 这里常常从database 或 redis 取数据
		bf := new(book.BookInfo)
		bf.BookId =  req.BookId
		bf.BookName = "天龙八部"
		fmt.Printf("makeGetBookInfoEndpoint bf:%+v\n",bf)
		return bf,nil
	}
}

func decodeRequest(_ context.Context,req interface{})(interface{},error){
	return req,nil
}

func encodeResponse(_ context.Context,rsp interface{})(interface{},error){
	return rsp,nil
}

func main(){
	//当前IP
	ip = getExternal()
	instance = ip+":56555"
	//内网地址即可 server 与 client 可以通信到就可以 127.0.0.1
	instance = "10.10.2.38:56555"
	value =  instance
	key = prefix+"/etcd001"

	options := etcdv3.ClientOptions{
		DialTimeout: time.Second*3,
		DialKeepAlive:time.Second*3,
	}

	//创建etcdv3 链接
	client,err := etcdv3.NewClient(ctx,[]string{etcdServer},options)
	if err != nil {
		panic(err)
	}
	fmt.Println("start........")


	//创建注册器
	registar :=  etcdv3.NewRegistrar(client,etcdv3.Service{
		Key: key,
		Value:value,
	},log.NewNopLogger())
	fmt.Printf("registar:%+v\n",registar)

	//注册器启动注册
	registar.Register()

	bookServer := new(BookServer)

	var edList endpoint.Endpoint
	edList = makeGetBookListEndpoint()
	//fmt.Printf("Endpoint response:%+v\n",edList)
	bookListHandler := grpc_transport.NewServer(
		edList,
		decodeRequest,
		encodeResponse,
	)
	bookServer.bookListHandler = bookListHandler
	fmt.Printf("bookListHandler:%+v\n",bookListHandler)

	var  edp endpoint.Endpoint
	edp = makeGetBookInfoEndpoint()
	fmt.Printf("Endpoint response:%+v\n",edp)
	bookInfoHandle := grpc_transport.NewServer(
		edp,
		decodeRequest,
		encodeResponse,
	)
	bookServer.bookInfoHandler = bookInfoHandle
	fmt.Printf("bookInfoHandler:%+v\n",bookInfoHandle)

	//serviceAddress = ip+":56555"
	ls,_ := net.Listen("tcp",serviceAddress)
	gs := grpc.NewServer(grpc.UnaryInterceptor(grpc_transport.Interceptor))

	fmt.Printf("gs:%+v\n",gs)
	book.RegisterBookServiceServer(gs,bookServer)
	fmt.Println("mid........")
	//http2 走这个  func (s *Server) newHTTP2Transport(c net.Conn, authInfo credentials.AuthInfo) transport.ServerTransport {
	gs.Serve(ls)
	fmt.Println("end........")
}