package main

import (
	"bufio"
	command "core/command/pb"
	"core/component/router"
	"fmt"
	"net"
	"time"
)

var (
	ch = make(chan []byte, 10)
	r  *router.Router
)

func runClientTcp(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}
	r = router.NewRouter()
	go send(conn)
	go recv(conn)
	add()
}

func recv(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf[:])
		if err != nil {
			return
		}
		msgRecv := &command.CSStartGameOnline{}
		r.UnMarshal(buf[:n], msgRecv)
		fmt.Println("------", msgRecv)
	}
}

func send(conn *net.TCPConn) {
	for {
		select {
		case data := <-ch:
			conn.Write(data)
			fmt.Println("========", string(data))
		}
	}
}

func add() {

	if r == nil {
		return
	}
	i := uint32(0)
	for {
		msgSend := &command.CSStartGameOnline{}
		msgSend.RoomId = 10000 + i
		data, err := r.Marshal(uint16(command.Command_CSStartGame), msgSend)
		if err != nil {
			fmt.Println(err)
			return
		}
		ch <- data
		time.Sleep(10 * time.Second)
		i++
	}
}
