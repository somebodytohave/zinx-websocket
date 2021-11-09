package znet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

//Message 消息
type Message struct {
	DataLen uint16 //消息的长度
	ID      uint16 //消息的ID 对应 AddRouter(1, &zrouter.PingRouter{})
	Data    []byte //消息的内容
	MsgType int    //消息类型websocket使用
}

//NewMsgPackage 创建一个Message消息包
func NewMsgPackage(ID uint16, msgType int, data []byte) *Message {
	return &Message{
		DataLen: uint16(len(data)),
		ID:      ID,
		Data:    data,
		MsgType: msgType,
	}
}

//NewBinaryMsgPackage 创建一个Message消息包
func NewBinaryMsgPackage(ID uint16, data []byte) *Message {
	return &Message{
		DataLen: uint16(len(data)),
		ID:      ID,
		Data:    data,
		MsgType: websocket.BinaryMessage,
	}
}

//获取消息类型
func (msg *Message) GetMsgType() int {
	return msg.MsgType
}

//GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() uint16 {
	return msg.DataLen
}

//GetMsgID 获取消息ID
func (msg *Message) GetMsgID() uint16 {
	return msg.ID
}

//GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len uint16) {
	msg.DataLen = len
}

//SetMsgID 设计消息ID
func (msg *Message) SetMsgID(msgID uint16) {
	msg.ID = msgID
}

//SetData 设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

//SetMsgType 设置消息类型 websocket
func (msg *Message) SetMsgType(msgType int) {
	msg.MsgType = msgType
}

//SetMsgType 设置消息类型 websocket
func (msg *Message) ToString() {
	fmt.Printf("msgID = %v, GetMsgType = %v, GetDataLen = %v, GetData = %v",
		msg.GetMsgID(), msg.GetMsgType(), msg.GetDataLen(), cast.ToString(msg.GetData()))
	fmt.Println()
}
