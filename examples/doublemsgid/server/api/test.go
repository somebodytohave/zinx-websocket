package api

import (
	"fmt"
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/ziface"
	"go.uber.org/zap"

	"github.com/sun-fight/zinx-websocket/znet"
)

//test 自定义路由
type TestRouter struct {
	znet.BaseRouter
}

//Handle
func (this *TestRouter) Handle(request ziface.IRequest) {
	global.Glog.Debug("Call TestRouter Handle")
	global.Glog.Debug("recv from client : ", zap.Any("msgid", request.GetMsgID()),
		zap.Any("data", string(request.GetData())))

	switch request.GetMsgID() {
	case 1001: //登录
		fmt.Println("主子命令-执行登录")
	case 1002: //退出登录
		fmt.Println("主子命令-执行退出登录")
	}
}
