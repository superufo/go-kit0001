package lib

import (
	"errors"
	"fmt"
)


/***
	相当于程序启动 Goroutine 而必须持有的令牌
 */
type GoTickets interface {
	//拿走一张票
	Take()
	//归还一张票
	Return()
	//票池是否被激活
	Active() bool
	//票的总数
	Total() uint32
	//剩余的票数
	Remainder() uint32
}

type myGotickets struct {
	// 票的总数
	total  int32
	//票的容器
	ticketCh  chan []struct{}
	//票池是否激活
	active bool
}

/*** myGotickets实现GoTickets 接口开始**/
func (gt  *myGotickets) Take(){
	<-gt.ticketCh
}

func (gt *myGotickets)Return(){
	gt.ticketCh <- struct{}{}
}

func (gt *myGotickets) Active() bool{
	return gt.active
}

func (gt *myGotickets)Total() uint32{
	return gt.total
} 

func (gt *myGotickets) Remainder() uint32{
	return uint32(len(gt.ticketCh))
}
/*** myGotickets实现GoTickets 接口结束**/

/*** myGotickets私有函数 **/
func (gt *myGotickets)init(total uint32) bool{
	if gt.active==false{
		return  false
	}

	if total == 0 {
		return false
	}

	ch := make(chan struct{}, total)
	n := int(total)
	for i:=0;i<n;i++ {
		ch <- struct{}{}
	}

	gt.ticketCh = ch
	gt.total = total
	gt.active = true

	return true
}

//驱动创建新的 ggoroutine 票池
//在Go中 一般以New 开头来初始化一类较复杂的架构体和接口
//依据面向接口的编程 一般不会直接将结构体返回
// 一般返回其实现的接口 这样程序的扩展性更好
//所有 接口的方法都是公用的，也只需要暴露几个被调用的方法即可
func NewGoTickets(total uint32) (GoTickets,error) {
	  gt := myGotickets{}
	  if !gt.init(total) {
	  	errMsg := fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
	  	return nil,errors.New(errMsg)
	  }

	  return &gt,nil
}




