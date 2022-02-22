package main

//func (tcpSrv *TcpServer) msgToOnline(msgId uint16, data []byte) {
//	msgRecv := &command.GatewayMsgToOnline{}
//	err := tcpSrv.router.UnMarshal(data, msgRecv)
//	if err != nil {
//		return
//	}
//	onlineStream, ok := tcpSrv.onlineStreams[msgRecv.RemoteAddr]
//	if !ok {
//		return
//	}
//	onlineStream.addSendClientMsgChan(msgRecv.Data)
//}
//
//func (tcpSrv *TcpServer) msgToGateway(msgId uint16, data []byte) {
//	msgRecv := &command.OnlineMsgToGateway{}
//	err := tcpSrv.router.UnMarshal(data, msgRecv)
//	if err != nil {
//		return
//	}
//	tcpConn, ok := tcpSrv.tcpConnects[msgRecv.RemoteAddr]
//	if !ok {
//		return
//	}
//	tcpConn.addMsgChannel(msgRecv.Data)
//}
