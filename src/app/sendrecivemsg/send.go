package sendrecivemsg

func SendByteToClient(byte []byte) {
	SendMsgChan <- byte
}
