package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sun-fight/zinx-websocket/examples/doublemsgid/server/api"
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/znet"
)

func main() {
	server := znet.NewServer()
	//比如 已有命令号 1001登录  1002退出登录
	// 1 = 主命令 = 1001/DoubleMsgID. 配置表zinx.json DoubleMsgID:1000
	mainCmd := 1001 / global.Object.DoubleMsgID
	//testRouter解析主命令包含所有 1xxx的子命令
	server.AddRouter(mainCmd, &api.TestRouter{})

	bindAddress := fmt.Sprintf("%s:%d", global.Object.Host, global.Object.TCPPort)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/ws", server.Serve)
	router.Run(bindAddress)
}
