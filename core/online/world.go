package main

import (
	command "core/command/pb"
	"core/component/logger"
	"sync"
)

var (
	world     *World
	worldOnce sync.Once
)

type World struct {
	players                   sync.Map
	fromOtherServerMsgMsgChan chan *command.ClientPlayerMsgData
}

func WorldGetMe() *World {
	worldOnce.Do(func() {
		world = &World{
			players:                   sync.Map{},
			fromOtherServerMsgMsgChan: make(chan *command.ClientPlayerMsgData, 1024),
		}
		go world.loop()
	})
	return world
}

//接收其他服转发过来的消息
func (w *World) sendFromOtherServerMsgChan(playerId uint64, data []byte) {
	w.fromOtherServerMsgMsgChan <- &command.ClientPlayerMsgData{PlayerId: playerId, Data: data}
}

func (w *World) loop() {
	for {
		select {
		case msg := <-w.fromOtherServerMsgMsgChan:
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
	logger.Infof("玩家(%v)进入 online", playerId)
	w.players.Store(playerId, player)
	logger.Infof("添加玩家(%v) onlineWorld", playerId)
}

func (w *World) delPlayer(playerId uint64) {
	w.players.Delete(playerId)
	logger.Infof("移除玩家(%v) onlineWorld", playerId)
}
