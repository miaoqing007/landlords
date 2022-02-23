package main

import (
	command "core/command/pb"
	"log"
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
			fromGatewayMsgChan: make(chan *command.ClientPlayerMsgData, 1024),
		}
		go world.loop()
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

func (w *World) addPlayer(playerId uint64, clientAddr, gatewayGRPCAddr string) {
	if _, ok := w.players.Load(playerId); ok {
		return
	}
	player := newPlayer(playerId, clientAddr, gatewayGRPCAddr)
	w.players.Store(playerId, player)
	log.Printf("添加玩家成功==", player)
}
