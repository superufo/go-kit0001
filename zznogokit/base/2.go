package main

import (
	"fmt"
	"os"
	"strings"
)

type UpperString string
func (s UpperString) String() string {
	return strings.ToUpper(string(s))
}

/*defer 语句延迟执行了一个匿名函数，因为这个匿名函数捕获了外部函数的
局部变量 v ，这种函数我们一般叫闭包。闭包对捕获的外部变量并不是传值方式
访问，而是以引用的方式访问*/
func main(){
	fmt.Fprintln(os.Stdout, UpperString("hello, world"))

	k := Inc()
	fmt.Printf("\nK: %d",k)

	for i := 0; i < 3; i++ {
		defer func(){ fmt.Printf( "1111 %d\n ",i) } ()
	}

	for i := 0; i < 3; i++ {
		i := i // 定义一个循环体内局部变量i
		defer func(){ fmt.Printf( "2222 %d\n ",i) } ()
	}

	for i := 0; i < 3; i++ {
		// 通过函数传入i
		// defer 语句会马上对调用参数求值
		defer func(i int){ fmt.Printf( "3333 %d\n ",i) } (i)
	}

}

func Inc() (v int) {
	defer func(){
		v++
		fmt.Printf("%d",v)
    } ()
	return 42
}



