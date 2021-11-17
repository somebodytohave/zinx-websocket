package api

import (
	"github.com/sun-fight/zinx-websocket/examples/protobuf/pb"
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/ziface"
	
	"github.com/sun-fight/zinx-websocket/znet"
	"google.golang.org/protobuf/proto"
)

type HeartRouter struct {
	znet.BaseRouter
}

// Handle
func (this *HeartRouter) Handle(request ziface.IRequest) {
	msg := pb.ReqHeart{}
	err := proto.Unmarshal(request.GetData(), &msg)
	if err != nil {
		global.Glog.Error(err.Error())
	}
	global.Glog.Debug(msg.String())

	marshal, err := proto.Marshal(&msg)
	if err != nil {
		global.Glog.Error(err.Error())
	}
	err = request.GetConnection().SendBinaryBuffMsg(request.GetMsgID(), marshal)
	if err != nil {
		global.Glog.Error(err.Error())
	}
}
