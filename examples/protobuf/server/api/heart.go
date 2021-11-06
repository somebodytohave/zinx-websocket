package api

import (
	"github.com/sun-fight/zinx-websocket/examples/protobuf/pb"
	"github.com/sun-fight/zinx-websocket/ziface"
	"github.com/sun-fight/zinx-websocket/zlog"
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
		zlog.Error(err)
	}
	zlog.Debug(msg.String())

	marshal, err := proto.Marshal(&msg)
	if err != nil {
		zlog.Error(err)
	}
	err = request.GetConnection().SendBinaryBuffMsg(request.GetMsgID(), marshal)
	if err != nil {
		zlog.Error(err)
	}
}
