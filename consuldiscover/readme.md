1>go get github.com/pborman/uuid   
  go get github.com/hashicorp/consul    

2> register.go    
func NewRegistrar(client Client, r *stdconsul.AgentServiceRegistration, logger log.Logger) *Registrar 
5个参数 注册中心consul的ip、端口，算术服务的本地ip和端口，日志记录工具     

3> server.go 增加健康检查方法 HealthCheck   
   并且需要在 中间件 ArithmeticService、loggingMiddleware、metricMiddleware添加实现   
   	//健康检查  
	HealthCheck() bool      
   // ArithmeticService实现HealthCheck   
    
参考：https://juejin.im/post/5c740a335188257c1e2c86a7   
 
D:\gopromod\zz.wt\consuldiscover\discover   
main.exe  -consul.host 47.112.111.171 -consul.port 8500  -service.host 110.235.246.150   -service.port 9550          
main.exe   -consul.host 47.112.111.171  -consul.port 8500 