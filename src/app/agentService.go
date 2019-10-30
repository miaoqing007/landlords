package main

import (
	"app/client_handler"
	"app/helper/stack"
	"app/misc/packet"
	"app/session"
	"bufio"
	"github.com/golang/glog"
	"net"
	"os"
)

func agentRun() {
	lestener, err := net.Listen("tcp", ServerHost+":"+ServerPort)
	if err != nil {
		glog.Info("listen error:", err)
		os.Exit(1)
	}
	defer lestener.Close()
	glog.Info("listening on " + ServerHost + ":" + ServerPort)
	for {
		conn, err := lestener.Accept()
		if err != nil {
			glog.Info("accept error:", err)
			os.Exit(1)
		}
		glog.Infof("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		sess := session.NewSession()
		go handleRequest(conn, sess)
		go handWriteResp(conn, sess)
	}
}

func handleRequest(conn net.Conn, sess *session.Session) {
	ip := conn.RemoteAddr().String()
	defer func() {
		glog.Info("disconnect:" + ip)
		sess.AddDieChan()
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			glog.Info("err=", err)
			return
		}
		reader := packet.Reader(b)
		c, err := reader.ReadS16()
		if err != nil {
			glog.Info("err=", err)
			return
		}
		bytes := executeHandler(c, sess, reader)
		for _, byt := range bytes {
			sess.AddSendChan(byt)
		}
	}
}

func executeHandler(code int16, sess *session.Session, reader *packet.Packet) [][]byte {
	defer stack.PrintRecoverFromPanic()
	handle := client_handler.Handlers[code]
	if handle == nil {
		return nil
	}
	retByte := handle(sess, reader)
	return retByte
}

func handWriteResp(conn net.Conn, sess *session.Session) {
	ch := make(chan []byte, 1)
	sess.EvaluationSendChan(ch)
	writer := bufio.NewWriter(conn)
	for {
		select {
		case msg := <-ch:
			writer.Write(msg)
			writer.Write([]byte("\n"))
			writer.Flush()
		case <-sess.Die:
			glog.Info("disconnect Write")
			return
		}
	}
}
