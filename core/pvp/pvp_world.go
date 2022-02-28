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
	players    sync.Map
	msgChannel chan *command.Online2PvpInfo
}

func WorldGetMe() *World {
	worldOnce.Do(func() {
		world = &World{
			msgChannel: make(chan *command.Online2PvpInfo, 1024),
		}
		go world.loop()
	})
	return world
}

func (w *World) getPlayer(playerId uint64) *Player {
	if value, ok := w.players.Load(playerId); ok {
		return value.(*Player)
	}
	return nil
}

func (w *World) addMsgChannel(msg *command.Online2PvpInfo) {
	w.msgChannel <- msg
}

func (w *World) addPlayer(playerId uint64, onlineRemoteAddr string) *Player {
	player := newPvpPlayer(playerId, onlineRemoteAddr)
	if _, ok := w.players.Load(playerId); !ok {
		w.players.Store(playerId, player)
		logger.Infof("玩家(%v)进入pvpWorld online-->pvp", playerId)
	}
	return player
}

func (w *World) delPlayer(playerId uint64) {
	w.players.Delete(playerId)
	logger.Infof("玩家(%v)移除pvpWorld", playerId)
}

func (w *World) loop() {
	for {
		select {
		case msg := <-w.msgChannel:
			player := w.getPlayer(msg.PlayerId)
			if player == nil {
				player = w.addPlayer(msg.PlayerId, msg.Addr)
			}
			player.addInnerMsgChannel(msg.Data)
		}
	}
}
