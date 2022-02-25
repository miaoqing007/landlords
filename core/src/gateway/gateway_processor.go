package main

import (
	command "core/command/pb"
	"core/component/logger"
)

func (tcpConn *TcpConn) clientInOnlineHandler(msgId uint16, data []byte) {
	msgRecv := &command.ClientInOnline_Online{}
	err := tcpConn.router.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	tcpConn.bindPlayerId(12345)

	msgRecv.PlayerId = tcpConn.playerId
	tcpConn.sendMsgToOnline(command.Command_ClientInOnline, msgRecv)
	logger.Infof("玩家(%v)进入online", msgRecv.PlayerId)
}

func (tcpConn *TcpConn) clientOutOnlineHandler(msgId uint16, data []byte) {
	msgRecv := &command.ClientOutOnline_Online{}
	err := tcpConn.router.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	logger.Infof("玩家(%v)进入online", msgRecv.PlayerId)
}
