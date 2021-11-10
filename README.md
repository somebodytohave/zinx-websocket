# 基于zinx框架二次开发(websocket版)

- tcp协议改为websocket
- 新增心跳检测功能

[zinx(TCP版本)](https://github.com/aceld/zinx)

[看云-《zinx框架教程-基于Golang的轻量级并发服务器》](https://www.kancloud.cn/aceld/zinx)

[简书-《zinx框架教程-基于Golang的轻量级并发服务器》](https://www.jianshu.com/p/23d07c0a28e5)

## 案例demo

### [ping-服务器与客户端的简单通信demo](https://github.com/sun-fight/zinx-websocket/tree/master/examples/ping)

### 一、快速上手

[代码来自examples->ping](https://github.com/sun-fight/zinx-websocket/tree/master/examples/ping)

### server端

基于zinx框架开发的服务器应用，主函数步骤比较精简，最多只需要3步即可。

```go
func main() {
//1 创建一个server句柄
server := znet.NewServer()

//2 配置路由
server.AddRouter(1, &api.PingRouter{})

//3 开启服务
bindAddress := fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
router := gin.Default()
router.GET("/", server.Serve)
router.Run(bindAddress)
}
```

其中(api.PingRouter)自定义路由及业务处理：
[代码跳转](https://github.com/sun-fight/zinx-websocket/blob/master/examples/ping/server/api/ping.go)

### client端

zinx的消息处理采用，`[MsgLength]|[MsgID]|[Data]`的封包格式
[代码跳转](https://github.com/sun-fight/zinx-websocket/blob/master/examples/ping/client/main.go)

### zinx配置文件

[详细配置文件说明与默认值](https://github.com/sun-fight/zinx-websocket/blob/master/utils/globalobj.go)

```json
{
  "Name": "zinx-websocket Demo",
  "Host": "127.0.0.1",
  "TcpPort": 8999,
  "MaxConn": 3,
  "WorkerPoolSize": 10,
  "LogDir": "./mylog",
  "LogFile": "zinx.log",
  "HeartbeatTime": 60
}
```

### I.服务器模块Server

```go
  func NewServer () ziface.IServer 
```

创建一个zinx服务器句柄，该句柄作为当前服务器应用程序的主枢纽，包括如下功能：

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
func (s *Server) AddRouter (msgId uint16, router ziface.IRouter) 
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