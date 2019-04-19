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

