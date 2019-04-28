package lib

import "time"

//RawReq  表示原生请求的结构
type RawReq struct {
	ID int64
	Req []byte
}

// 定义 原生响应的结构
type RawResp struct{
	ID int64
	Resp []byte
	Err error
	Elapse time.Duration
}

// Retcode 表示结果代码的类型
type RetCode int

type CallResult struct {
	ID int64               //ID。
	Req RawReq             //原生请求
	Resp RawResp           //原生响应
	Code RetCode           //响应代码
	Msg string             //结果成因的简述
	Elapse time.Duration   //耗时
}

// 保留 1-10000 给在和承受方使用
const (
	RET_CODE_SUCCESS             RetCode= 0        // 成功
	RET_CODE_WARNING_CALL_TIMEOUT       = 1001    // 调用超时警告
	RET_CODE_ERROR_CALL                 = 2001    // 调用错误
	RET_CODE_ERROR_RESPONSE				= 2002    // 响应内容错误
	RET_CODE_ERROR_CALEE				= 2003    // 被调用方（被测试软件）内部错误
	RET_CODE_FATAL_CALL					= 3001    // 调用过程发生的致命错误
)

// s根据结果代码返回响应的文字解释
func GetRetCodePlain(code RetCode) string {
	var codePlain string
	switch code {
	case  RET_CODE_SUCCESS:
		codePlain = "Successs"
	case  RET_CODE_WARNING_CALL_TIMEOUT:
		codePlain = " Call Timeout Warning"
	case  RET_CODE_ERROR_CALL:
		codePlain = "Call Error"
	case  RET_CODE_ERROR_RESPONSE:
		codePlain = "Response Error"
	case  RET_CODE_ERROR_CALEE:
		codePlain = "Callee Error"
	case  RET_CODE_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown Result Code"
	}

	return codePlain
}

// 声明荷载发生器状态的常量
const   (
	//代表原始
	STATUS_ORIGINAL uint32 = 0
	//正在启动
	STATUS_STARTING uint32 = 1
	//已经启动
	STATUS_STARTED  uint32 = 2
	//正在停止
	STATUS_STOPPING uint32 =3
	//已经停止
	STATUS_STOPPED  uint32 =4
)

//定义荷载发生器的接口
type Generator interface {
	//启动负载发生器 结果值代表是否已经成功启动
	Start()      bool
	//停止荷载发生器  结果值代表是否已经成功停止
	Stop()       bool
	//获取 状态
	Status()     bool
	// 获取调用的计数 每次启动会重置该计数
	CallCount()  int64
}