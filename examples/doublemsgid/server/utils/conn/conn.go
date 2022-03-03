package utils

import (
	"github.com/sun-fight/zinx-websocket/global"
	"github.com/sun-fight/zinx-websocket/ziface"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func SendBinaryBuffMsg(res proto.Message, request ziface.IRequest) (err error) {
	return SendBinaryBuffMsgByMsgID(res, request, 0)
}

func SendBinaryBuffMsgByMsgID(res proto.Message, request ziface.IRequest, msgID uint16) (err error) {
	var u16MsgID uint16
	if msgID > 0 {
		u16MsgID = uint16(msgID)
	} else {
		u16MsgID = request.GetMsgID()
	}
	bytes, err := proto.Marshal(res)
	if err != nil {
		global.Glog.Error("解析返回数据,", zap.Uint16("msgID", u16MsgID), zap.Any("data", res), zap.Error(err))
		return
	}
	err = request.GetConnection().SendBinaryBuffMsg(u16MsgID, bytes)
	if err != nil {
		global.Glog.Error("发送消息,", zap.Uint16("msgID", u16MsgID), zap.Any("data", res), zap.Error(err))
		return
	}
	global.Glog.Info("发送消息", zap.Int64("connID = ", request.GetConnection().GetConnID()),
		zap.Uint16("msgID", u16MsgID), zap.Any("data", res))
	return
}
