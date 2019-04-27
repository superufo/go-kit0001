package flow

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
	"os"
	"strconv"
)

func Register(consulHost,consulPort,svcHost,svcPort string,logger log.Logger)(register sd.Registrar){
	//创建 Consul 客户端连接  此处 {}表示 consul.Client 接口的实例化
	var client consul.Client
	{
		 consulCfg := api.DefaultConfig()
         consulCfg.Address = consulHost+":"+consulPort
		 consulClient,err := api.NewClient(consulCfg)

		 if err!=nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		 }

		 client = consul.NewClient(consulClient)
	}

	//设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP: "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "11s",
		Timeout: "30s",
		Notes:"Consul check service health status.",
	}

	port,_ := strconv.Atoi(svcPort)
	reg := api.AgentServiceRegistration{
		ID : "arithmetic"+uuid.New(),
		Name :  "arithmetic",
		Address: svcHost,
		Port: port,
		Tags:[]string{"arithmetic", "mike"},
		Check:&check,
	}

	//执行注册
    register = consul.NewRegistrar(client,&reg,logger)
	return
}