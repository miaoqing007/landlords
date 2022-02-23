package main

import (
	"core/component/router"
	"fmt"
)

type Player struct {
	innerRouter     *router.Router
	isInitFinish    bool
	playerChanMsg   chan []byte
	closeChannel    chan bool
	gatewayGRPCAddr string
	clientAddr      string
	//User          *UserManager
}

func newPlayer(playerId uint64, clientAddr, gatewayGRPCAddr string) *Player {
	player := &Player{
		playerChanMsg:   make(chan []byte, 64),
		closeChannel:    make(chan bool, 0),
		gatewayGRPCAddr: gatewayGRPCAddr,
		clientAddr:      clientAddr,
	}
	go player.loop()
	return player
}

func (player *Player) loop() {
	if !player.isInitFinish {
		player.init()
		player.isInitFinish = true
	}
	for {
		select {
		case data := <-player.playerChanMsg:
			player.onMessage(data)
		case <-player.closeChannel:
			return
		}
	}
}

func (player *Player) addPlayerChanMsg(data []byte) {
	player.playerChanMsg <- data
}

func (player *Player) onMessage(data []byte) {
	player.innerRouter.Route(data)
}

func (player *Player) init() {
	player.initBase()
	player.regMsgHandler()
	WorldGetMe()
}

func (player *Player) initBase() error {
	player.innerRouter = router.NewRouter()
	//userManger, err := NewUserManager(account, password)
	//if err != nil {
	//	return errors.New("")
	//}
	//player.User = userManger
	return nil
}

func (player *Player) sendMSg(msgId uint16, msg interface{}) {
	data, err := player.innerRouter.Marshal(msgId, msg)
	if err != nil {
		return
	}
	gatewayStream := tcpServer().getOnlineStream(player.gatewayGRPCAddr)
	if gatewayStream == nil {
		return
	}
	gatewayStream.addToGatewayMsg(data, player.clientAddr)
	fmt.Println(data)
}
