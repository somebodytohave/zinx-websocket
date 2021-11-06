package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sun-fight/zinx-websocket/examples/protobuf/server/api"
	"github.com/sun-fight/zinx-websocket/utils"
	"github.com/sun-fight/zinx-websocket/znet"
)

func main() {
	server := znet.NewServer()
	server.AddRouter(1, &api.HeartRouter{})

	bindAddress := fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/", server.Serve)
	router.Run(bindAddress)
}
