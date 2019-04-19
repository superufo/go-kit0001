`httputil.ReverseProxy` 编写反向代理最重要的就是实现自己的 `Director`

> Director must be a function which modifies the request into a new request to be sent using Transport.
> Its response is then copied back to the original client unmodified.
> Director must not access the provided Request after returning.

Director` 是一个函数，它接受一个请求作为参数，然后对其进行修改。修改后的请求会实际发送给服务器端，
因此我们编写自己的 `Director` 函数，每次把请求的 Scheme 和 Host 修改成某个后端服务器的地址，
就能实现负载均衡的效果（其实上面的正向代理也可以通过相同的方法实现）

curl 请求

![](https://github.com/superufo/go-kit0001/blob/master/zznogokit/img/rpn0001.png)

结果：

![](https://github.com/superufo/go-kit0001/blob/master/zznogokit/img/rpn0002.png)

![](https://github.com/superufo/go-kit0001/blob/master/zznogokit/img/rpn0003.png)







