package main

import (
	"crypto/md5"
	"io"
	"strconv"
	"time"
	"fmt"
)

func main(){
	//这样处理后，因为我们限定了修改只能使用POST，当GET方式请求时就拒绝响应

	CsrfToken()
}


/***
***  csrf
***  防范在非GET请求中增加伪随机数；服务器先将生成的伪随机数保存到redis，将生成的伪随机数增加到表单字段，
*** 传输给浏览器，浏览器提交后，对比此伪随机数是否在redis ,验证是否是正确的客户端
*** 此方法也可以防止表单重复提交
*** https://www.kancloud.cn/kancloud/web-application-with-golang/44166  防止表单重复提交
*** https://www.kancloud.cn/kancloud/web-application-with-golang/44194 csrf
 */
func CsrfToken(){
	crutime := time.Now().Unix()
	fmt.Println(crutime)
	h:= md5.New()
	fmt.Println(h)
	io.WriteString(h,strconv.FormatInt(crutime,10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Println(token)
}

