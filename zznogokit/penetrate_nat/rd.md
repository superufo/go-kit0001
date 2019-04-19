**MultiWriter 的使用**

```
file, err := os.Create("tmp.txt")
if err != nil {
    panic(err)
}
defer file.Close()
writers := []io.Writer{
	file,
	os.Stdout,
}
writer := io.MultiWriter(writers...)
writer.Write([]byte("Go语言中文网"))
```

这段程序执行后在生成 tmp.txt 文件，同时在文件和屏幕中都输出：`Go语言中文网`。这和 Unix 中的 tee 命令类似。

**Copy 函数**

```
func Copy(dst Writer, src Reader) (written int64, err error)
```

Copy 将 src 复制到 dst，直到在 src 上到达 EOF 或发生错误。它返回复制的字节数.

成功的 Copy 返回 err == nil，而非 err == EOF

```
io.Copy(os.Stdout, strings.NewReader("Go语言中文网"))
```

直接将内容输出（写入 Stdout 中）

我们甚至可以这么做：

```
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	io.Copy(os.Stdout, os.Stdin)
	fmt.Println("Got EOF -- bye")
}
```

执行：`echo "Hello, World" | go run main.go`

**strings对象**      --<https://www.jianshu.com/p/1675c3fddabd>

主要做字符串操作：Replace   ContainsAny ContainsAny  Fields FieldsFunc Split  HasPrefix HasSuffix

Index  IndexAny  IndexFunc  IndexRune LastIndex LastIndexAny LastIndexFunc Join Repeat

```
func Fields(s string) []string
func FieldsFunc(s string, f func(rune) bool) []string
Fields 用一个或多个连续的空格分隔字符串s
```

```
func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }
func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) }
func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
通过 sep 进行分割，返回[]string
```

```
func Join(a []string, sep string) string
假如没有这个库函数，我们自己实现一个，我们会这么实现：
```

**log**

```
log.SetFlags 定义输出格式
const (
   Ldate         = 1 << iota     // 日期:  2009/01/23
   Ltime                         // 时间:  01:23:23
   Lmicroseconds                 // 微秒:  01:23:23.123123.
   Llongfile                     // 路径+文件名+行号: /a/b/c/d.go:23
   Lshortfile                    // 文件名+行号:   d.go:23
   LUTC                          // 使用标准的UTC时间格式
   LstdFlags     = Ldate | Ltime // 默认
)
```

**使用**

```
nb -listen 1997 2017 -log D:/nb
nb -tran 1997 192.168.1.2:338 -log D:/nb
nb -slave 127.0.0.1:3389 8.8.8.8:1997 -log D:/nb
```

假设有外网主机`123.123.123.123:1997`和`123.123.123.123:2017`端口开放。

内网主机`192.168.1.2:3389`需要转发到外网。首先在外网主机执行

```
nb -listen 1997 2017
```

作用是开辟两个用于监听内网打隧道的连接端口和其他应用客户端连接的端口。

接着内网主机执行

```
nb -slave 127.0.0.1:3389 123.123.123.123:1997
```

作用是内网主机主动连接外网主机打通隧道。

然后其他客户端（例如本例子中的3389远程桌面客户端）连接`123.123.123.123:2017`，就等同于连接到了内网主机的`192.168.1.2:3389`上。



参考：<https://www.cnblogs.com/jkko123/p/7218685.html>

​           <https://www.zhihu.com/people/changwei1997/activities>

​         <https://github.com/cw1997/NATBypass>