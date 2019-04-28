package lib

import "time"

// 不知道调用接口的api 定义此接口 可能是http 或者 tcp 等等
// 组件包括 请求的生成操作和 响应的检查操作
// 叫做调用器接口
type  Caller interface {
	// 构建请求
	BuildReq() RawReq
	//调用
	Call(req []byte,timeoutNs time.Duration) ([]byte,error)
	// 检查响应
	CheckResp(rawReq RawReq,rawResp RawResp) *CallResult
}
