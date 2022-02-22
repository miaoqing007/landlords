package main

import (
	command "core/command/pb"
	"sync"
)

var (
	world     *World
	worldOnce sync.Once
)

type World struct {
	players            sync.Map
	fromGatewayMsgChan chan *command.ClientPlayerMsgData
}

func WorldGetMe() *World {
	worldOnce.Do(func() {
		world = &World{
			players:            sync.Map{},
			fromGatewayMsgChan: make(chan *command.ClientPlayerMsgData, 64),
		}
		world.loop()
	})
	return world
}

func (w *World) sendFromGatewayMsgChan(msg *command.ClientPlayerMsgData) {
	w.fromGatewayMsgChan <- msg
}

func (w *World) loop() {
	for {
		select {
		case msg := <-w.fromGatewayMsgChan:
			player := WorldGetMe().getPlayer(msg.PlayerId)
			if player != nil {
				player.addPlayerChanMsg(msg.Data)
			}
		}
	}
}

func (w *World) getPlayer(playerId uint64) *Player {
	player, ok := w.players.Load(playerId)
	if !ok {
		return nil
	}
	return player.(*Player)
}
