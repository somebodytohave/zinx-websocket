# 基于Zinx框架二次开发(websocket版)

- tcp协议改为websocket
- 新增心跳检测功能

# Zinx简介

[Zinx 是一个基于Golang的轻量级并发服务器框架](https://github.com/aceld/zinx)

[看云-《Zinx框架教程-基于Golang的轻量级并发服务器》](https://www.kancloud.cn/aceld/zinx)

[简书-《Zinx框架教程-基于Golang的轻量级并发服务器》](https://www.jianshu.com/p/23d07c0a28e5)

## 案例demo
### [Ping-简单的服务器与客户端的交互](https://github.com/sun-fight/zinx-sun/tree/master/examples/ping)

### 一、快速上手

[代码来自examples->ping](https://github.com/sun-fight/zinx-sun/tree/master/examples/ping)
### server端
基于Zinx框架开发的服务器应用，主函数步骤比较精简，最多只需要3步即可。

```go
func main() {
//1 创建一个server句柄
server := znet.NewServer()

//2 配置路由
server.AddRouter(0, &PingRouter{})

//3 开启服务
bindAddress := fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
router := gin.Default()
router.GET("/", server.Serve)
router.Run(bindAddress)
}
```

其中自定义路由及业务配置方式如下：

```go
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

```

### client端

zinx-sun的消息处理采用，`[MsgLength]|[MsgID]|[Data]`的封包格式
[对比修改了拆包的逻辑](https://github.com/gorilla/websocket/blob/master/examples/echo/client.go)

```go
package main

import (
	"flag"
	"fmt"
	"github.com/sun-fight/zinx-sun/znet"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8999", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			//获取 io读取流
			msgType, ioReader, err := c.NextReader()
			if err != nil {
				fmt.Println("get read reader error ", err)
				return
			}
			//读取客户端的Msg head
			dataPack := znet.NewDataPack()
			headData := make([]byte, dataPack.GetHeadLen())
			if _, err := io.ReadFull(ioReader, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
			//拆包，得到msgID 和 datalen 放在msg中
			msg, err := dataPack.Unpack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}
			msg.SetMsgType(msgType)

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(ioReader, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)
			//打印数据
			msg.ToString()
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			msgPackage := znet.NewBinaryMsgPackage(1, []byte(t.String()))
			pack, err := znet.NewDataPack().Pack(msgPackage)
			if err != nil {
				log.Println("write:", err)
				return
			}
			err = c.WriteMessage(websocket.BinaryMessage, pack)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
```

### Zinx配置文件

```json
{
  "Name": "zinx-sun ping Demo",
  "Host": "127.0.0.1",
  "TcpPort": 8999,
  "MaxConn": 3,
  "WorkerPoolSize": 10,
  "LogDir": "./mylog",
  "LogFile": "zinx.log",
  "HeartbeatTime": 30
}
```

`Name`:服务器应用名称

`Host`:服务器IP

`TcpPort`:服务器监听端口

`MaxConn`:允许的客户端链接最大数量

`WorkerPoolSize`:工作任务池最大工作Goroutine数量

`LogDir`: 日志文件夹

`LogFile`: 日志文件名称(如果不提供，则日志信息打印到Stderr)

`HeartbeatTime`: 心跳检测时间默认30秒

### I.服务器模块Server

```go
  func NewServer () ziface.IServer 
```

创建一个Zinx服务器句柄，该句柄作为当前服务器应用程序的主枢纽，包括如下功能：

#### 1)开启服务

```go
  func (s *Server) Start(c *gin.Context)
```

#### 2)停止服务

```go
  func (s *Server) Stop()
```

#### 3)运行服务

```go
  func (s *Server) Serve(c *gin.Context)
```

#### 4)注册路由

```go
func (s *Server) AddRouter (msgId uint32, router ziface.IRouter) 
```

#### 5)注册链接创建Hook函数

```go
func (s *Server) SetOnConnStart(hookFunc func (ziface.IConnection))
```

#### 6)注册链接销毁Hook函数

```go
func (s *Server) SetOnConnStop(hookFunc func (ziface.IConnection))
```

### II.路由模块

```go
//实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseRouter)PreHandle(req ziface.IRequest){}
func (br *BaseRouter)Handle(req ziface.IRequest){}
func (br *BaseRouter)PostHandle(req ziface.IRequest){}
```

### III.链接模块

#### 1)获取原始的socket TCPConn

```go
  func (c *Connection) GetTCPConnection() *net.TCPConn 
```

#### 2)获取链接ID

```go
  func (c *Connection) GetConnID() uint32 
```

#### 3)获取远程客户端地址信息

```go
  func (c *Connection) RemoteAddr() net.Addr 
```

#### 4)发送消息

```go
func (c *Connection) SendMsg(msgID uint32, msgType int, data []byte) error
func (c *Connection) SendBuffMsg(msgID uint32, msgType int, data []byte) error
//默认二进制消息
func (c *Connection) SendBinaryMsg(msgID uint32, data []byte) error
func (c *Connection) SendBinaryBuffMsg(msgID uint32, data []byte) error
```

#### 5)链接属性

```go
//设置链接属性
func (c *Connection) SetProperty(key string, value interface{})

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error)

//移除链接属性
func (c *Connection) RemoveProperty(key string) 
```

---