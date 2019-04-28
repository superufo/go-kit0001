package main

import (
"flag"
"fmt"
"os"
)

// 实际中应该用更好的变量名
var (
	h bool
	v bool
	V bool
	t bool
	T bool
	q *bool

	s string
	p string
	c string
	g string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&V, "V", false, "show version and configure options then exit")

	flag.BoolVar(&t, "t", false, "test configuration and exit")
	flag.BoolVar(&T, "T", false, "test configuration, dump it and exit")

	// 另一种绑定方式
	q = flag.Bool("q", false, "suppress non-error messages during configuration testing")

	// 注意 `signal`。默认是 -s string，有了 `signal` 之后，变为 -s signal
	flag.StringVar(&s, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")
	flag.StringVar(&p, "p", "/usr/local/nginx/", "set `prefix` path")
	flag.StringVar(&c, "c", "conf/nginx.conf", "set configuration `file`")
	flag.StringVar(&g, "g", "conf/nginx.conf", "set global `directives` out of configuration file")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
	}else{
		fmt.Println("v:",v)
		fmt.Println("t:",t)
		fmt.Println("s:",s)
	}


	married := flag.Bool("married", false, "Are you married?")
	age := flag.Int("age", 22, "How old are you?")
	name := flag.String("name", "", "What your name?")

	var address string
	//flag.StringVar这样的函数第一个参数换成了变量地址，后面的参数和flag.String是一样的。
	flag.StringVar(&address, "address", "GuangZhou", "Where is your address?")

	flag.Parse() //解析输入的参数

	fmt.Println("输出的参数married的值是:", *married)//不加*号的话,输出的是内存地址
	fmt.Println("输出的参数age的值是:", *age)
	fmt.Println("输出的参数name的值是:", *name)
	fmt.Println("输出的参数address的值是:", address)
}

func usage() {
	fmt.Fprintf(os.Stderr, `nginx version: nginx/1.10.0
Usage: nginx [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
`)
	flag.PrintDefaults()
}
