package main

import "sync"

type Room struct {
	roomPlayers sync.Map //map[uint64]*Player
}
