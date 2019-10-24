package sendrecivemsg

func SendMsgToClient(msg []byte) {
	SendMsgChan <- msg
}
