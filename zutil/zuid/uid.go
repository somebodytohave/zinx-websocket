package zuid

import (
	"github.com/bwmarrin/snowflake"
	"github.com/sun-fight/zinx-websocket/global"
)

const (
	NodeZinx = 0
)

var _node *snowflake.Node

func Init() {
	snowflake.NodeBits = 0
	snowflake.StepBits = 22
	if snowflake.NodeBits+snowflake.StepBits > 22 {
		panic("snowflake NodeBits add StepBits must less than 22")
	}
	var err error
	newNode, err := snowflake.NewNode(NodeZinx)
	if err != nil {
		global.Glog.Fatal("new snowflake node " + err.Error())
	}
	_node = newNode
}
func Gen64() int64 {
	return _node.Generate().Int64()
}

func Gen() snowflake.ID {
	return _node.Generate()
}
