package znet

import (
	"fmt"
	"github.com/sun-fight/zinx-websocket/global"
	"time"
)

//定时检测心跳包
func (c *Connection) heartBeatChecker() {
	if global.Object.HeartbeatTime == 0 {
		return
	}
	var (
		timer *time.Timer
	)

	timer = time.NewTimer((global.Object.HeartbeatTime) * time.Second)

	for {
		select {
		case <-timer.C:
			if !c.IsAlive() {
				c.Stop()
				//心跳检测失败，结束连接
				fmt.Println("连接已关闭 或者 太久没有心跳")
				return
			}
			timer.Reset(global.Object.HeartbeatTime * time.Second)
		case <-c.ctx.Done():
			timer.Stop()
			fmt.Println("连接已关闭")
			return
		}
	}

}

//IsAlive 检测心跳
func (c *Connection) IsAlive() bool {
	var (
		now = time.Now()
	)
	c.Lock()
	defer c.Unlock()
	if c.isClosed || now.Sub(c.lastHeartBeatTime) >
		global.Object.HeartbeatTime*time.Second {
		return false
	}
	return true

}

//KeepAlive 更新心跳
func (c *Connection) KeepAlive() {
	var (
		now = time.Now()
	)
	c.Lock()
	defer c.Unlock()

	c.lastHeartBeatTime = now
}
