package api

import (
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/ziface"
	"github.com/sun-fight/zinx-websocket/znet"
	"go.uber.org/zap"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {

	global.Glog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	global.Glog.Debug("recv from client : ", zap.Any("msgid", request.GetMsgID()),
		zap.Any("data", request.GetData()))

	err := request.GetConnection().SendBinaryBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		global.Glog.Error(err.Error())
	}
}
