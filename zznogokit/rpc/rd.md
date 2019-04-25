reverserrpc      反向rpc     

kvstore   带watch功能的 rpc





rpcpb02 来自   4.1.2 更安全的RPC接口

接口规范中针对客户端新增加了HelloServiceClient类型，该类型也必须满足
HelloServiceInterface接口，这样客户端用户就可以直接通过接口对应的方法调用
RPC函数。同时提供了一个 DialHelloService 方法，直接拨号 HelloService 服务

