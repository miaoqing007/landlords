package main

import (
	command "core/command/pb"
	"core/component/logger"
	"core/component/router"
)

type Player struct {
	router               *router.Router
	onlineGRPCRemoteAddr string
	roomId               uint64
	playerId             uint64
	innerMsgChannel      chan []byte //玩家消息通道
	closeChannel         chan bool   //退出pvp信号
}

func newPvpPlayer(playerId uint64, onlineRemoteAddr string) *Player {
	p := &Player{
		router:               router.NewRouter(),
		playerId:             playerId,
		onlineGRPCRemoteAddr: onlineRemoteAddr,
		innerMsgChannel:      make(chan []byte, 128),
		closeChannel:         make(chan bool, 0),
	}
	p.registerMsgHandler()
	go p.loop()
	return p
}

func (player *Player) addInnerMsgChannel(data []byte) {
	player.innerMsgChannel <- data
}

func (player *Player) loop() {
	defer func() {
		WorldGetMe().delPlayer(player.playerId)
	}()
	for {
		select {
		case data := <-player.innerMsgChannel:
			player.onMessage(data)
		case <-player.closeChannel:
			logger.Infof("玩家(%v)退出pvp", player.playerId)
			return
		}
	}
}

func (player *Player) onMessage(data []byte) {
	player.router.RouterOnlinePvpMsg(data)
}

func (player *Player) sendMsgToOnlineClient(cmd command.Command, msg interface{}) {
	data, err := player.router.Marshal(uint16(cmd), msg)
	if err != nil {
		return
	}
	stream := gSrvGetMe().getOnlineStream(player.onlineGRPCRemoteAddr)
	if stream == nil {
		return
	}
	stream.addMsgToOnine(player.playerId, data)
}
