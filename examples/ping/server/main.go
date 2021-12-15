package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sun-fight/zinx-websocket/examples/ping/server/api"
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/znet"
)

func main() {
	server := znet.NewServer()
	server.AddRouter(1, &api.PingRouter{})

	bindAddress := fmt.Sprintf("%s:%d", global.Object.Host, global.Object.TCPPort)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/ws", server.Serve)
	router.Run(bindAddress)
}
