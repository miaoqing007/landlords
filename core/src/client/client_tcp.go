package main

import (
	"bufio"
	command "core/command/pb"
	"core/component/logger"
	"core/component/router"
	"core/config"
	"fmt"
	"net"
	"time"
)

var (
	ch = make(chan []byte, 10)
	r  *router.Router
)

func runClientTcp() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+config.GatewayTCPPort)
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
	defer func() {
		conn.Close()
	}()

	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf[:])
		if err != nil {
			return
		}
		msgRecv := &command.CSStartGame_Online{}
		r.UnMarshal(buf[:n], msgRecv)
		fmt.Println("------", msgRecv)
	}
}

func send(conn *net.TCPConn) {
	for {
		select {
		case data := <-ch:
			if _, err := conn.Write(data); err != nil {
				return
			}
		}
	}
}

func add() {

	if r == nil {
		return
	}

	var (
		msgSend interface{}
		//cmd     command.Command
		//i       = uint32(0)
	)
	msgSend = &command.ClientInOnline_Online{}
	ch <- sendMsg(command.Command_ClientInOnline, msgSend)
	time.Sleep(3 * time.Second)
	//for {
	//	if i%2 == 0 {
	//		msgSend = &command.CSJoinPvpPool_Pvp{}
	//		msgSend.(*command.CSJoinPvpPool_Pvp).PlayerId = 80000 + uint64(i)
	//		cmd = command.Command_CSJoinPvpPool
	//	} else {
	//		msgSend = &command.CSStartGame_Online{}
	//		msgSend.(*command.CSStartGame_Online).RoomId = 10000 + i
	//		cmd = command.Command_CSStartGame
	//	}
	//	ch <- sendMsg(cmd, msgSend)
	//	time.Sleep(10 * time.Second)
	//	i++
	//	if i == 6 {
	//		break
	//	}
	//}
	msgSend = &command.ClientOutOnline_Online{}
	ch <- sendMsg(command.Command_ClientOutOnline, msgSend)
	time.Sleep(100 * time.Second)
}

func sendMsg(cmd command.Command, msg interface{}) []byte {
	fmt.Println("---==-=-=-", msg)
	data, err := r.Marshal(uint16(cmd), msg)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return data
}
