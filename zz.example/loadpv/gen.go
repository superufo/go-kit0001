package loadpv

import (
	"bytes"
	"context"
	"errors"
	"math"
	"sync/atomic"
	"time"
	"fmt"

	"loadpv/helper/log"
	"loadpv/lib"
)

/***
*并发量 = 单个在和的响应超时时间 / 载荷的发送间隔时间
* 单个在和的响应超时时间 即timeoutNs 字段的值表示的时间.
*一旦操作载荷的耗时达到了响应超时时间,
*该载荷就会被判定为未被成功响应,在这之后,
* 载荷发生器不会等待该载荷的响应,假设响应超时时间是5秒,
* 在响应超时时间为5秒的设定下,如果每隔1秒向被测试软件发送一个载荷.
* 那么这个并发量就是5, 而如果每隔1MS发送一个载荷,那么该并发量就是5000.
* 这里所说的发送间隔时间可以由在和发生器的lps字段的值计算得出
* cocurrency = timeoutNs / (1e9/lps ) +1
* le9 代表的是 1秒对应的纳秒数, 加1 代表 周期之初,向被测试软件发送的那个载荷.
***/

 //日志记录器
 var logger = log.DLogger()

 //载荷发生器的实现模型
type myGenerator struct{
	//调用器
	caller        lib.Caller
	//处理超时的时间 单位纳秒
	timeoutNS     time.Duration
	//每秒载荷量 每一秒发送的请求数
	lps           uint32
	//负载持续时间 单位纳秒
	durationNS    time.Duration
	// 载荷并发量。
	concurrency uint32
	//Goroutine 票池
	tickets       lib.GoTickets
	//上下文
	ctx           context.Context
	//取消函数
	cancelFunc    context.CancelFunc
	//调用计数
	callCount     int64
	//状态
	status        uint32
	//调用结果通道
	resultCh      chan *lib.CallResult
}

 //**** myGenerator 实现Generator 接口开始**//
// 启动在和发生器
func (gen *myGenerator) Start() bool{
	logger.Infoln("Starting load generator....")

	//CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
	//参数addr指向的被操作值与参数old的值是否相等。当相等时new代表的新值替换掉原先的旧值
	if !atomic.CompareAndSwapInt32(&gen.status,lib.STATUS_ORIGINAL,lib.STATUS_STARTING){
		if !atomic.CompareAndSwapInt32(&gen.status,lib.STATUS_STOPPED,lib.STATUS_STARTING){
			return false
		}
	}

	//设定节流阀  1e9 一秒代表的纳秒数  time.Duration 转成int64
    var throttle <-chan time.Time
	if gen.lps > 0 {
		// 一秒对应纳秒内发送的请求 得出一个请求发送需要间隔多少纳秒
		interval := time.Duration(1e9/gen.lps)
		logger.Infoln("Setting throttle (%v)... ",interval)
		//等待interval 纳秒会 返回通道值给throttle
		throttle = time.Tick(interval)
	}

    //初始化上下文和取消函数
    gen.ctx,gen.cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)

    //初始化调用计数
    gen.callCount = 0

    //设置状态为已经启动
    atomic.StoreUint32(&gen.status,lib.STATUS_STARTED)

    go func() {
    	//生成并发送载荷
    	logger.Infoln("Generating loads....")
    	gen.genLoad(throttle)
		logger.Infof("Stopped. (call count: %d)", gen.callCount)
	}()

	return true
}

 func (gen *myGenerator) Stop() bool{
 	if !atomic.CompareAndSwapInt32(&gen.status,lib.STATUS_STARTED,lib.STATUS_STOPPING){
		return false
	}
 	gen.cancelFunc()

 	for{
		if atomic.LoadUint32(&gen.status)== lib.STATUS_STOPPED{
			break
		}
		time.Sleep(time.Microsecond)
	}

 	return true
 }

 func (gen *myGenerator)Status() uint32{
 	//atomic原子操作类都是对地址操作
 	return atomic.LoadUint32(&gen.status)
 }

 func (gen *myGenerator)CallCount() int64{
 	return atomic.LoadInt64(&gen.callCount)
 }
//*** myGenerator 实现Generator接口结束**//


//*** myGenerator私有函数开始 **//
// 初始化载荷发生器。
func (gen *myGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("Initializing the load generator...")
	// 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
	var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets

	buf.WriteString(fmt.Sprintf("Done. (concurrency=%d)", gen.concurrency))
	logger.Infoln(buf.String())
	return nil
}

// asyncSend 会异步地调用承受方接口。
 func  (gen *myGenerator)asyncCall()() {
	 gen.tickets.Take()
	 go func() {
		 defer func() {
			 if p := recover(); p != nil {
			 	//此种语法结构不甚明白
				 err, ok := interface{}(p).(error)
				 var errMsg string
				 if ok {
					 errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
				 } else {
					 errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
				 }
				 logger.Errorln(errMsg)
				 result := &lib.CallResult{
					 ID:   -1,
					 Code: lib.RET_CODE_FATAL_CALL,
					 Msg:  errMsg}
				 gen.sendResult(result)
			 }
			 gen.tickets.Return()
		 }()

		 rawReq := gen.caller.BuildReq()
		 // 调用状态：0-未调用或调用中；1-调用完成；2-调用超时。
		 var callStatus uint32
		 //超时的处理方法
		 timer := time.AfterFunc(gen.timeoutNS, func() {
			 if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				 return
			 }
			 result := &lib.CallResult{
				 ID:     rawReq.ID,
				 Req:    rawReq,
				 Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
				 Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
				 Elapse: gen.timeoutNS,
			 }
			 gen.sendResult(result)
		 })

		 rawResp := gen.callOne(&rawReq)
		 //callStatus：0-未调用或调用中；1-调用完成
		 if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			 return
		 }

		 //Stop prevents the Timer from firing.
		 timer.Stop()
		 var result *lib.CallResult
		 if rawResp.Err != nil {
			 result = &lib.CallResult{
				 ID:     rawResp.ID,
				 Req:    rawReq,
				 Code:   lib.RET_CODE_ERROR_CALL,
				 Msg:    rawResp.Err.Error(),
				 Elapse: rawResp.Elapse}
		 } else {
			 result = gen.caller.CheckResp(rawReq, *rawResp)
			 result.Elapse = rawResp.Elapse
		 }
		 gen.sendResult(result)
	 }()
 }

// callOne 会向载荷承受方发起一次调用。
func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
	atomic.AddInt64(&gen.callCount, 1)
	if rawReq == nil {
		return &lib.RawResp{ID: -1, Err: errors.New("Invalid Raw Request.")}
	}
	//现在时间精确到纳秒
	start := time.Now().UnixNano()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
	end := time.Now().UnixNano()
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime}
	} else {
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Resp:   resp,
			Elapse: elapsedTime}
	}
	return &rawResp
}

// sendResult 用于发送调用结果。
func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
	if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}
	select {
	case gen.resultCh <- result:
		return true
	default:
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

// printIgnoredResult 打印被忽略的结果。
func (gen *myGenerator) printIgnoredResult(result *lib.CallResult, cause string) {
	resultMsg := fmt.Sprintf(
		"ID=%d, Code=%d, Msg=%s, Elapse=%v",
		result.ID, result.Code, result.Msg, result.Elapse)
	logger.Warnf("Ignored result: %s. (cause: %s)\n", resultMsg, cause)
}

func (gen *myGenerator)prepareToStop(ctxError error){
	logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
	atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
	logger.Infof("Closing result channel...")
	close(gen.resultCh)
	atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
}

// genLoad 会产生载荷并向承受方发送。
 func (gen *myGenerator)genLoad(throttle chan time.Time){
	for {
		select{
		case <- gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}

		gen.asyncCall()
		if gen.lps>0{
			select{
			case <-throttle:
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())\
				return
			}
		}
	}
}
//*** myGenerator私有函数结束 **//

// NewGenerator 会新建一个载荷发生器 驱动载荷发生器。
func NewGenerator(pset ParamSet) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}


