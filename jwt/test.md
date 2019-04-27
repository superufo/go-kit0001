```

D:\gopromod\zz.wt\jwt (master -> origin)
auth.exe  -consul.host 47.112.111.171 -consul.port 8500  -service.host 127.0.0.1   -service.port 9550     --zipkin.url http://47.112.111.171:9411/zipkin/api/v2/spans


测试：http://127.0.0.1:9550/login  {"name":"mike","pwd":"123456"}


```



![](https://github.com/superufo/go-kit0001/blob/master/img/jwt01.png)

![](https://github.com/superufo/go-kit0001/blob/master/img/jwt02.png)
