package main

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

type info struct {
	name string
	id int
}

func main()  {
	var a = "shijie  世界是老道那个任命"
	str2bytes(a)

	v := info{"Nan", 33}
	fmt.Printf("%v\n", v)
	fmt.Printf("%+v\n", v)

	var array3 = []int{9, 10, 11, 12}
	fmt.Printf("array3--- type:%T \n", array3)

	//创建数组切片，并仅初始化其中的部分元素，数组切片的len将根据初始化的元素确定
	var array6 = []string{4: "Smith", 2: "Alice"}
	fmt.Printf("array6--- type:%T \n", array6)
	rangeObjPrint(array6)

	var aa = []float64{4, 2, 5, 7, 2, 1, 88, 1}
	var f= SortFloat64FastV2(aa)
	for i,k :=range  f {
		fmt.Printf("key:%d  value:%d \n", i, k)
	}

	var m=SortFloat64FastV1(aa)
	for ii,kk:=range m{
		fmt.Printf("key:%d value%d \n",ii,kk )
	}


}

func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		p[i] = c
	}
	fmt.Println("p:%+v\n",p)

	//var c = []int{1, 2, 3}
	//var e = c[0:2:cap(c)]
	//fmt.Println("e %#v\n",reflect.SliceHeader(e))


	var a []int
	a = append(a, 1) // 追加1个元素
	a = append(a, 1, 2, 3) // 追加多个元素, 手写解包方式
	a = append(a, []int{1,2,3}...) // 追加一个切片, 切片需要解包

	return p
}


//输出整型数组切片
func rangeIntPrint(array []int) {
	for i, v := range array {
		fmt.Printf("index:%d  value:%d\n", i, v)
	}
}

//输出字符串数组切片
func rangeObjPrint(array []string) {
	for i, v := range array {
		fmt.Printf("index:%d  value:%s\n", i, v)
	}
}

func SortFloat64FastV1(a []float64) []int {
	// 强制类型转换
	var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(
		a):cap(a)]
	// 以int方式给float64排序
	sort.Ints(b)

	return b
}
func SortFloat64FastV2(a []float64) []int  {
	// 通过 reflect.SliceHeader 更新切片头部信息实现转换
	var c []int
	aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
	*cHdr = *aHdr
	// 以int方式给float64排序
	sort.Ints(c)

	for _, cc := range c {
		fmt.Printf(" cc value:%d\n", cc)
	}

	fmt.Printf("len=%d cap=%d slice=%v\n",len(c),cap(c),c)


	return c
}






