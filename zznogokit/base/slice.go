package main

import (
	"fmt"
)

func main(){
	s:= []string{"aaaa","ssss","ssssadss","vcxvx","fgrassas","4654sfsf","vsfaeq14"}
	fmt.Println( s[4] )
	fmt.Println( s[2:2] )
	fmt.Println( s[2:3] )
	fmt.Println( s[3:] )

	s1 := []int{1, 2, 3, 4, 5} //短操作符声明 len为4，cap为4
	s2 := make([]int, 2, 4)    //make语法声明 ，len为2，cap为4
	s2 = []int{5, 6}

	s3 := append(s2, 7) //append一个元素
	fmt.Println(s3, s2) //[5 6 7] [5 6]

	s4 := append(s2, s1...) //append  一个切片所有的元素
	fmt.Println(s4)         //[5 6 3 4]

	//return
	copy(s1, s2)    //  复制，用s2的元素填充s1里去，改变原slice,覆盖对应的key
	fmt.Println(s1) //[5 6 3 4]

	s1[0], s1[1] = 1, 2
	copy(s2, s1)
	fmt.Println(s2) //[1 2] 目标slice len不够时，只填满len

	s5 := s1[1:4]
	s6 := s5[0:4] //不会报错，因为cap为4，从底层取得最后一位

	fmt.Println(s5, s6, cap(s6)) //[2 3 4] [2 3 4 5] 4

	//删除第三个元素
	s7 := append(s1[:1], s1[3:]...)
	fmt.Println(s7) //[1 4 5]


}
