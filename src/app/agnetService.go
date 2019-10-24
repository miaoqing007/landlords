package main

import (
	"app/sendrecivemsg"
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func agentRun() {
	flag.Parse()
	var l net.Listener
	var err error
	l, err = net.Listen("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("listening on " + *host + ":" + *port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			os.Exit(1)
		}
		fmt.Printf("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
		go sendrecivemsg.RecieveMsgFromClient()
	}
}

func handleRequest(conn net.Conn) {
	ip := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnect:" + ip)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)
	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		//select {
		//case msg := <-sendrecivemsg.SendMsgChan:
		//	writer.Write(msg)
		//	writer.Flush()
		//}
		sendrecivemsg.ReciveMsgChan <- b
	}
	fmt.Println("Closed Server!")
}
