package main

//
//func (conn *TcpConn) msgToOnline(msgId uint16, data []byte) {
//	msgRecv := &command.ClientPlayerMsgData{}
//	err := conn.router.UnMarshal(data, msgRecv)
//	if err != nil {
//		return
//	}
//	onineStream := tcpServer().getOnlineStream(conn.onlineStreamAddr)
//	if onineStream == nil {
//		return
//	}
//	onineStream.addToClientMsg(msgRecv)
//}
//
//func (conn *TcpConn) msgToGateway(msgId uint16, data []byte) {
//	msgRecv := &command.ClientPlayerMsgData{}
//	err := conn.router.UnMarshal(data, msgRecv)
//	if err != nil {
//		return
//	}
//	//tcpConn, ok := tcpServer().tcpConnects[conn.ipAddr]
//	//if !ok {
//	//	return
//	//}
//	conn.addMsgChannel(msgRecv.Data)
//}
