package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/server/zrouter"
	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/znet"
	"github.com/gin-gonic/gin"
)

func main() {
	server := znet.NewServer()
	server.AddRouter(1, &zrouter.PingRouter{})

	bindAddress := fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/", server.Serve)
	r.Run(bindAddress)

}
