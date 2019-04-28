package main

import  (
  "fmt"
)
func main() {
	done := make(chan []int, 10) // 带 10 个缓存
	// 开N个后台打印线程
	for i := 0; i < cap(done); i++ {
		go func(){
			fmt.Println("你好, 世界")
			done <- []int{1,3,5,6,6}
		}()
	}

	// 等待N个后台线程完成  不可能连续打印 3个  "你好, 人类"
	for i := 0; i < cap(done); i++ {
		fmt.Println("你好, 人类")
		var tt= <-done
		fmt.Println("tt:",tt)
	}
}
