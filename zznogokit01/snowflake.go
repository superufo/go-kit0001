package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
)

/**
** 产生  1 bit unused + 41 timestample + 10 bit node Id + 12 bit sequence id
** timestamp ， datacenter_id ， worker_id 和 sequence_id 这四个字段中， timestamp 和 sequence_id 是由程序在运行期生成的。
** 但 datacenter_id 和 worker_id 需要我们在部署阶段就能够获取得到，并且一旦程序启动之后，就是不可更改的了（想想，如果可以随意更改，可能被不慎修改，造成最终生成的id有冲突）
**/
func main() {
	n, err := snowflake.NewNode(1)

	if err != nil {
		println(err)
		os.Exit(1)
	}

	for i := 0; i < 3; i++ {
		id := n.Generate()
		fmt.Println("id", id)
		fmt.Println("node: ", id.Node(), "step: ", id.Step(), "time: ", id.Time(), "\n", )
	}
}

/*********
// Epoch is set to the twitter snowflake epoch of Nov 04 201
0 01:42:54 UTC
// You may customize this to set a different epoch for your
application.
Epoch int64 = 1288834974657
// Number of bits to use for Node
// Remember, you have a total 22 bits to share between Node/
Step
NodeBits uint8 = 10
// Number of bits to use for Step
// Remember, you have a total 22 bits to share between Node/
Step
StepBits uint8 = 12
Epoch 就是本节开头讲的起始时间， NodeBits 指的是机器编号的位
长， StepBits 指的是自增序列的位长
 *******/


