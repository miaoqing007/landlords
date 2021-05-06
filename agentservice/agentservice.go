package agentservice

import (
	"github.com/golang/glog"
	"landlords/client_handler"
	"landlords/enmu"
	"landlords/helper/stack"
	"landlords/misc/packet"
	"landlords/session"
	"net"
	"os"
)

func AgentRun() {
	lestener, err := net.Listen("tcp", ":"+enmu.ServerPort)
	if err != nil {
		glog.Info("listen error:", err)
		os.Exit(1)
	}
	defer lestener.Close()
	glog.Info("listening on " + enmu.ServerHost + ":" + enmu.ServerPort)
	for {
		conn, err := lestener.Accept()
		if err != nil {
			glog.Info("accept error:", err)
			os.Exit(1)
		}
		glog.Infof("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	in := make(chan []byte, 16)
	sess := session.NewSession(in)
	//registry.Register()
	defer func() {
		glog.Info("disconnect:" + conn.RemoteAddr().String())
		sess.OffLine()
		conn.Close()
	}()
	go func() {
		for {
			select {
			case msg := <-in:
				conn.Write(msg)
			}
		}
	}()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		c, data := packet.UnPacket(buf[:n])
		in <- executeHandler(c, sess, data)
	}
}

//////执行方法
func executeHandler(code int16, sess *session.Session, data []byte) []byte {
	defer stack.PrintRecoverFromPanic()
	handle := client_handler.Handlers[code]
	if handle == nil {
		return nil
	}
	return packet.Pack(handle(sess, data))
}
