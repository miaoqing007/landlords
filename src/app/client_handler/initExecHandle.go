package client_handler

import (
	"app/helper/stack"
	"app/misc/packet"
	"app/session"
)

func InitHandle() {
	session.ExecuteHandler = executeHandler
}

func executeHandler(code int16, sess *session.Session, reader *packet.Packet) [][]byte {
	defer stack.PrintRecoverFromPanic()
	handle := Handlers[code]
	if handle == nil {
		return nil
	}
	retByte := handle(sess, reader)
	return retByte
}
