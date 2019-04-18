#### Endpoint

        // Endpoint is the fundamental building block of servers and clients.
        // It represents a single RPC method.
        type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
         
        // Nop is an endpoint that does nothing and returns a nil error.
        // Useful for tests.
        func Nop(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil }
         
        // Middleware is a chainable behavior modifier for endpoints.
        type Middleware func(Endpoint) Endpoint
         
        // Chain is a helper function for composing middlewares. Requests will
        // traverse them in the order they're declared. That is, the first middleware
        // is treated as the outermost middleware.
        func Chain(outer Middleware, others ...Middleware) Middleware {
            return func(next Endpoint) Endpoint {
                for i := len(others) - 1; i >= 0; i-- { // reverse
                    next = others[i](next)
                }
                return outer(next)
            }
        }

####常见main.go模式代码 （http）
     package main

         import  (
            "context"
            "flag"
            "fmt"
            "golang.org/x/time/rate"
            "net/http"
            "os"
            "os/signal"
          
            "syscall"
            "time"
          
            "github.com/go-kit/kit/log"
          
            kitprometheus "github.com/go-kit/kit/metrics/prometheus"
            stdprometheus "github.com/prometheus/client_golang/prometheus"
            "register/flow"
         )
      
        func main(){
         var (
                consulHost  = flag.String("consul.host", "", "consul ip address")
                consulPort  = flag.String("consul.port", "", "consul port")
                serviceHost = flag.String("service.host", "", "service ip address")
                servicePort = flag.String("service.port", "", "service port")
         )
        flag.Parse()
        
        ctx := context.Background()
        errChan := make(chan error)
      
        var logger log.Logger
        {
            logger = log.NewLogfmtLogger(os.Stderr)
            logger = log.With(logger, "ts", log.DefaultTimestampUTC)
            logger = log.With(logger, "caller", log.DefaultCaller)
        }
      
        fieldKeys := []string{"method"}
        requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
            Namespace:"mike",
            Subsystem:"arithmetic_service",
            Name:      "request_count",
            Help:      "Number of requests received.",
        },fieldKeys)
      
        requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
            Namespace: "mike",
            Subsystem: "arithemetic_service",
            Name:      "request_latency",
            Help:      "Total duration of requests in microseconds.",
        }, fieldKeys)
      
        var svc  flow.Service
        svc = flow.ArithmeticService{}
      
        //包裹日志 进行装饰
        svc = flow.LogginMiddleware(logger)(svc)
      
        //Api监控 进行装饰
        svc = flow.Metrics(requestCount, requestLatency)(svc)
      
        //传入请求 ，传出响应（）， 调用service 处理业务逻辑
        endpoint := flow.MakeArithmeticEndpoint(svc)


​      
​        //包裹 限流器 判断流量是否溢出
​        // add ratelimit,refill every second,set capacity 3
​        //ratebucket := ratelimit.NewBucket(time.Second*1, 3)
​        //endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)
​        // 一秒3次的流量  怎么传进来  就怎么传出去  中间做了限流的异常判断
​        ratelimitter := rate.NewLimiter(rate.Every(time.Second*1),100)
​         endpoint = flow.NewTokenBucketLimitterWithBuildIn(ratelimitter)(endpoint)
​      
         /*************** 健康检查 ***********************/
         //创建健康检查的Endpoint，未增加限流
         healthEndpoint := flow.MakeHealthCheckEndpoint(svc)
      
        //把算术运算Endpoint和健康检查Endpoint封装至ArithmeticEndpoints
        endpts := flow.SetsEndpoints{
            ArithmeticEndpoint:  endpoint,
            HealthCheckEndpoint: healthEndpoint,
        }
      
        //创建http.Handler 处理传输协议 （http,tcp grpc等） 处理返回的格式
        r := flow.MakeHttpHandler(ctx, endpts, logger)
      
        //创建注册对象
        registar := flow.Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)
      
        go func() {
            fmt.Println("Http Server start at port:" + *servicePort)
            //启动前执行注册
            registar.Register()
            handler := r
            errChan <- http.ListenAndServe(":"+*servicePort, handler)
        }()
      
        go func() {
            c := make(chan os.Signal, 1)
            signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
            errChan <- fmt.Errorf("%s", <-c)
        }()
      
        error := <-errChan
      
        //服务退出取消注册
        registar.Deregister()
        fmt.Println(error)
 	}

#### main.go grpc模式代码：
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


​     
​    	//创建注册器
​    	registar :=  etcdv3.NewRegistrar(client,etcdv3.Service{
​    		Key: key,
​    		Value:value,
​    	},log.NewNopLogger())
​    	fmt.Printf("registar:%+v\n",registar)
​     
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
    	gs.Serve(ls)
    	fmt.Println("end........")
    }

### 1.go-kit 支持的协议

![img](https://img-blog.csdnimg.cn/20190418142926177.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

 nats 协议的介绍: <https://www.cnblogs.com/liang1101/p/6641082.html>  支持同步消息和异步消息

AMQP协议:            <http://www.amqp.org/>

### 2.http request response 

### ![img](https://img-blog.csdnimg.cn/20190418131137630.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/20190418131544553.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

 D:\gopro\pkg\mod\github.com\go-kit\kit@v0.8.0\transport\http

![img](https://img-blog.csdnimg.cn/20190418131855382.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/20190418132242826.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/20190418132847366.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/2019041813294959.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/20190418133117611.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

写一个服务的流程： encode Http reponse :

![img](https://img-blog.csdnimg.cn/20190418140550195.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/20190418141221798.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

### 3.grpc  request response 

![img](https://img-blog.csdnimg.cn/20190418145311220.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

![img](https://img-blog.csdnimg.cn/20190418145136718.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

```
bookListHandler := grpc_transport.NewServer(
		edList,
		decodeRequest,
		encodeResponse,
	)
 
bookServer.bookListHandler = bookListHandler
fmt.Printf("bookListHandler:%+v\n",bookListHandler)
```

```
type BookServer struct{
	bookListHandler grpc_transport.Handler
	bookInfoHandler grpc_transport.Handler
}
```

github.com\go-kit\kit@v0.8.0\transport\grpc 中定义

```
type Handler interface {
	ServeGRPC(ctx oldcontext.Context, request interface{}) (oldcontext.Context, interface{}, error)
}
```

```
 grpc_transport.NewServer(
   edList,
   decodeRequest,
   encodeResponse,
)   
结果的调用基于以下代码 (github.com\go-kit\kit@v0.8.0\transport\grpc)
```

```
// Server wraps an endpoint and implements grpc.Handler.
type Server struct {
	e         endpoint.Endpoint
	dec       DecodeRequestFunc
	enc       EncodeResponseFunc
	before    []ServerRequestFunc
	after     []ServerResponseFunc
	finalizer []ServerFinalizerFunc
	logger    log.Logger
}
 
// NewServer constructs a new server, which implements wraps the provided
// endpoint and implements the Handler interface. Consumers should write
// bindings that adapt the concrete gRPC methods from their compiled protobuf
// definitions to individual handlers. Request and response objects are from the
// caller business domain, not gRPC request and reply types.
func NewServer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *Server {
	s := &Server{
		e:      e,
		dec:    dec,
		enc:    enc,
		logger: log.NewNopLogger(),
	}
	for _, option := range options {
		option(s)
	}
	return s
}
```

edList 方法 与  方法decodeRequest 绑定代码 （也就是  e 与  dec 绑定 ） 

在  func (s Server) ServeGRPC(ctx oldcontext.Context, req interface{})

(retctx oldcontext.Context, resp interface{}, err error) 中实现  （其中s 是 s := &Server 的s）：

![img](https://img-blog.csdnimg.cn/20190418150659338.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3N1cGVyX3Vmbw==,size_16,color_FFFFFF,t_70)

源码：

```
// ServeGRPC implements the Handler interface.
func (s Server) ServeGRPC(ctx oldcontext.Context, req interface{}) (retctx oldcontext.Context, resp interface{}, err error) {
	// Retrieve gRPC metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
 
	if len(s.finalizer) > 0 {
		defer func() {
			for _, f := range s.finalizer {
				f(ctx, err)
			}
		}()
	}
 
	for _, f := range s.before {
		ctx = f(ctx, md)
	}
 
	var (
		request  interface{}
		response interface{}
		grpcResp interface{}
	)
 
	request, err = s.dec(ctx, req)
	if err != nil {
		s.logger.Log("err", err)
		return ctx, nil, err
	}
 
	response, err = s.e(ctx, request)
	if err != nil {
		s.logger.Log("err", err)
		return ctx, nil, err
	}
 
	var mdHeader, mdTrailer metadata.MD
	for _, f := range s.after {
		ctx = f(ctx, &mdHeader, &mdTrailer)
	}
 
	grpcResp, err = s.enc(ctx, response)
	if err != nil {
		s.logger.Log("err", err)
		return ctx, nil, err
	}
 
	if len(mdHeader) > 0 {
		if err = grpc.SendHeader(ctx, mdHeader); err != nil {
			s.logger.Log("err", err)
			return ctx, nil, err
		}
	}
 
	if len(mdTrailer) > 0 {
		if err = grpc.SetTrailer(ctx, mdTrailer); err != nil {
			s.logger.Log("err", err)
			return ctx, nil, err
		}
	}
 
	return ctx, grpcResp, nil
}
```



来自 superufo:

https://blog.csdn.net/super_ufo/article/details/89372098



md编辑工具：typora 推荐使用： 

https://typora.io/



下载：

git clone http://superufo:××××@github.com/superufo/go-kit0001.git  zz.wt