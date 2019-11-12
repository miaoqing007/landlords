package client_handler

import (
	"app/helper/stack"
	"app/misc/packet"
	"app/session"
	"github.com/golang/glog"
)

func InitHandle() {
	session.ExecuteHandler = executeHandler
	glog.Info("初始化handle完成")
}

//执行方法
func executeHandler(code int16, sess *session.Session, reader *packet.Packet) [][]byte {
	defer stack.PrintRecoverFromPanic()
	handle := Handlers[code]
	if handle == nil {
		return nil
	}
	retByte := handle(sess, reader)
	return retByte
}
