package api

import (
	"fmt"
	"github.com/sun-fight/zinx-websocket/ziface"
	"github.com/sun-fight/zinx-websocket/zlog"
	"github.com/sun-fight/zinx-websocket/znet"
)

//test 自定义路由
type TestRouter struct {
	znet.BaseRouter
}

//Handle
func (this *TestRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call TestRouter Handle")
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	switch request.GetMsgID() {
	case 1001: //登录
		fmt.Println("主子命令-执行登录")
	case 1002: //退出登录
		fmt.Println("主子命令-执行退出登录")
	}
}
