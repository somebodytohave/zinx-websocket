package znet

import (
	"fmt"
	"github.com/aceld/zinx/utils"
	"time"
)

//定时检测心跳包
func (c *Connection) heartBeatChecker() {
	var (
		timer *time.Timer
	)

	timer = time.NewTimer(time.Duration(utils.GlobalObject.HeartbeatTime) * time.Second)

	for {
		select {
		case <-timer.C:
			if !c.IsAlive() {
				c.Stop()
				//心跳检测失败，结束连接
				fmt.Println("连接已关闭 或者 太久没有心跳")
				return
			}
			timer.Reset(time.Duration(utils.GlobalObject.HeartbeatTime) * time.Second)
		case <-c.ctx.Done():
			timer.Stop()
			fmt.Println("连接已关闭")
			return
		}
	}

}

//检测心跳
func (c *Connection) IsAlive() bool {
	var (
		now = time.Now()
	)
	c.Lock()
	defer c.Unlock()
	if c.isClosed || now.Sub(c.lastHeartBeatTime) >
		time.Duration(utils.GlobalObject.HeartbeatTime)*time.Second {
		return false
	}
	return true

}

//更新心跳
func (c *Connection) KeepAlive() {
	var (
		now = time.Now()
	)
	c.Lock()
	defer c.Unlock()

	c.lastHeartBeatTime = now
}
