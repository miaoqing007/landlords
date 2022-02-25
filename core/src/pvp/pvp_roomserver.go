package main

import "sync"

type RoomSrv struct {
	rooms sync.Map //map[uint64]*Room
}
