package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/znet"
)

func main() {
	server := znet.NewServer()

	global.InitGormMysql()
	global.InitRedis()
	fmt.Println(global.Mysql)
	fmt.Println(global.Redis)

	bindAddress := fmt.Sprintf("%s:%d", global.Object.Host, global.Object.TCPPort)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/ws", server.Serve)
	router.Run(bindAddress)
}
