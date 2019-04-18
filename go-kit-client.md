### **go-kit client**  

------

**1.Http  服务发现的客户端**          

main.go    consuldiscover/discover/main.go

```go
     ........................

	ctx := context.Background()
                              
	//创建Endpoint  client 表示consul 客户端
	discoverEndpoint :=  flow.MakeDiscoverEndpoint(ctx, client, logger)

	//创建传输层
	r := flow.MakeHttpHandler(discoverEndpoint)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "9001")
		errc <- http.ListenAndServe(":9001", r)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}
```

discoverEndpoint :=  flow.MakeDiscoverEndpoint(ctx, client, logger)    原型(consuldiscover/discover/endpoints.go)：    

```go
func MakeDiscoverEndpoint(ctx context.Context, client consul.Client, logger log.Logger) endpoint.Endpoint 

// MakeDiscoverEndpoint 使用consul.Client创建服务发现Endpoint
// 为了方便这里默认了一些参数
func MakeDiscoverEndpoint(ctx context.Context, client consul.Client, logger log.Logger) endpoint.Endpoint {
	serviceName := "arithmetic"
	tags := []string{"arithmetic", "mike"}
	passingOnly := true
	duration := 500 * time.Millisecond

	//基于consul客户端、服务名称、服务标签等信息，
	// 创建consul的连接实例，
	// 可实时查询服务实例的状态信息
	instancer := consul.NewInstancer(client, logger, serviceName, tags, passingOnly)

	//针对calculate接口创建sd.Factory
	factory := arithmeticFactory(ctx, "POST", "calculate")

	//使用consul连接实例（发现服务系统）、factory创建sd.Factory
	endpointer := sd.NewEndpointer(instancer, factory, logger)

	//创建RoundRibbon负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	//为负载均衡器增加重试功能，同时该对象为endpoint.Endpoint
	retry := lb.Retry(1, duration, balancer)

	return retry
}
```

```go
//针对calculate接口创建sd.Factory
factory := arithmeticFactory(ctx, "POST", "calculate")
```

```go
func arithmeticFactory(_ context.Context,method,path string) sd.Factory{
   return func(instance string)(endpoint endpoint.Endpoint,closer io.Closer,err error) {
      if !strings.HasPrefix(instance, "http" ) {
        instance = "http://" + instance
      }

      tgt ,err := url.Parse(instance)
      if err!=nil {
         return nil,nil,err
      }

      tgt.Path = path

      var (
         enc kithttp.EncodeRequestFunc
         dec kithttp.DecodeResponseFunc
      )

      enc , dec = encodeArithmeticRequest, decodeArithmeticReponse

      return kithttp.NewClient(method,tgt,enc,dec).Endpoint(),nil,nil
   }
}
```

```
sd.Factory 定义 在 github.com\go-kit\kit@v0.8.0\sd\factory.go 中，定义如下：
```

```
type Factory func(instance string) (endpoint.Endpoint, io.Closer, error)
```

![1555583463245](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\1555583463245.png)

github.com\go-kit\kit@v0.8.0\transport\http\client.go 下面的Endpoint 方法：

```
// Endpoint returns a usable endpoint that invokes the remote endpoint.
func (c Client) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx, cancel := context.WithCancel(ctx)

		var (
			resp *http.Response
			err  error
		)
		if c.finalizer != nil {
			defer func() {
				if resp != nil {
					ctx = context.WithValue(ctx, ContextKeyResponseHeaders, resp.Header)
					ctx = context.WithValue(ctx, ContextKeyResponseSize, resp.ContentLength)
				}
				for _, f := range c.finalizer {
					f(ctx, err)
				}
			}()
		}

		req, err := http.NewRequest(c.method, c.tgt.String(), nil)
		if err != nil {
			cancel()
			return nil, err
		}

		if err = c.enc(ctx, req, request); err != nil {
			cancel()
			return nil, err
		}

		for _, f := range c.before {
			ctx = f(ctx, req)
		}

		resp, err = c.client.Do(req.WithContext(ctx))

		if err != nil {
			cancel()
			return nil, err
		}

		// If we expect a buffered stream, we don't cancel the context when the endpoint returns.
		// Instead, we should call the cancel func when closing the response body.
		if c.bufferedStream {
			resp.Body = bodyWithCancel{ReadCloser: resp.Body, cancel: cancel}
		} else {
			defer resp.Body.Close()
			defer cancel()
		}

		for _, f := range c.after {
			ctx = f(ctx, resp)
		}

		response, err := c.dec(ctx, resp)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
```

![1555583650121](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\1555583650121.png)



 标准库： net\http\request.go  NewRequest 定义如下：

```go
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	if method == "" {
		method = "GET"
	}
	if !validMethod(method) {
		return nil, fmt.Errorf("net/http: invalid method %q", method)
	}
	u, err := parseURL(url) // Just url.Parse (url is shadowed for godoc).
	if err != nil {
		return nil, err
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	// The host's colon:port should be normalized. See Issue 14836.
	u.Host = removeEmptyPort(u.Host)
	req := &Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(Header),
		Body:       rc,
		Host:       u.Host,
	}
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
```

![1555583790033](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\1555583790033.png)

encode request : 如下

![1555583953192](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\1555583953192.png)

decode response: 如下

​         ![1555584054111](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\1555584054111.png)



------

**2.grpc 服务发现的客户端**      