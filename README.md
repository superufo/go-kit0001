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