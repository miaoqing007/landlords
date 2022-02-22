package main

import command "core/command/pb"

type OnlinePvpStream struct {
	streams map[string]command.OnlinePvp_OnlinePvpStreamServer
}

func newOnlinePvpStream() *OnlinePvpStream {
	op := new(OnlinePvpStream)
	op.streams = make(map[string]command.OnlinePvp_OnlinePvpStreamServer)
	return op
}

func (op *OnlinePvpStream) addStream(addr string, streamServer command.OnlinePvp_OnlinePvpStreamServer) {
	if _, ok := op.streams[addr]; !ok {
		return
	}
	op.streams[addr] = streamServer

}
