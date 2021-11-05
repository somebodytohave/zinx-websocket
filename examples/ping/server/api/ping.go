package api

import (
	"github.com/sun-fight/zinx-sun/ziface"
	"github.com/sun-fight/zinx-sun/zlog"
	"github.com/sun-fight/zinx-sun/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {

	zlog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBinaryBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		zlog.Error(err)
	}
}
