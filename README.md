# go-kit0001
go-kit0001

原盘结构 D:\gopromod\zz.wt   
GOPATH   D:\gopro   
GOBIN    D:\gopro\bin
GOROOT   D:\Go\   
GOCACHE  D:\gopro\gocache
GoLand   D:\Program Files\JetBrains\GoLand 2018.3.5\bin     

交叉编译： 

```
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```

GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
GOARCH：目标平台的体系架构（386、amd64、arm）

非微服务的代码包：
zznogokit  主要是 advanced-go-programming-book.pdf 部分grpc 代码  
https://github.com/chai2010/advanced-go-programming-book/tree/master/examples             
 

zz.example                                
zznogokit01    

除此上述文件件都是微服务的代码      