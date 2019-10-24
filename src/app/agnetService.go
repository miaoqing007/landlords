package main

import (
	"app/session"
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func agentRun() {
	flag.Parse()
	lestener, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(1)
	}
	defer lestener.Close()
	fmt.Println("listening on " + *host + ":" + *port)
	for {
		conn, err := lestener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			os.Exit(1)
		}
		fmt.Printf("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		sess := session.NewSession()
		go handleRequest(conn, sess)
		go handWriteResp(conn, sess)
		go sess.ListenRecieveChan()
	}
}

func handleRequest(conn net.Conn, sess *session.Session) {
	ip := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnect:" + ip)
		sess.AddDieChan()
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		sess.AddRecieveChan(b)
	}
}

func handWriteResp(conn net.Conn, sess *session.Session) {
	ch := make(chan []byte, 1)
	sess.EvaluationSendChan(ch)
	for {
		select {
		case msg := <-ch:
			writer := bufio.NewWriter(conn)
			writer.Write(msg)
			writer.Write([]byte("\n"))
			writer.Flush()
		case <-sess.Die:

			return
		}
	}
}
