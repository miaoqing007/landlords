package sendrecivemsg

import (
	"app/clienthandle"
	"app/misc/packet"
	"app/session"
)

func RecieveMsgFromClient() {
	for msg := range ReciveMsgChan {
		reader := packet.Reader(msg)
		c, err := reader.ReadS16()
		if err != nil {
			return
		}
		executeHandler(c, &session.Session{}, reader)
	}
}

func executeHandler(code int16, sess *session.Session, reader *packet.Packet) [][]byte {
	//defer helper.PrintPanicStack()

	handle := clienthandle.Handlers[code]
	if handle == nil {
		return nil
	}
	//t := time.Now().UnixNano()
	// handle request
	retByte := handle(sess, reader)
	//testPack := client_handler.HandlersTime[code]
	//testPack.Num++
	//testPack.TotalTime = testPack.TotalTime + float64(time.Now().UnixNano()-t)/float64(time.Millisecond)
	return retByte
}
