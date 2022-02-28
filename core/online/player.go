package main

import (
	command "core/command/pb"
	"core/component/logger"
	"core/component/router"
)

type Player struct {
	innerRouter     *router.Router
	isInitFinish    bool
	playerChanMsg   chan []byte
	closeChannel    chan bool
	gatewayGRPCAddr string
	clientAddr      string
	playerId        uint64
	//User          *UserManager
}

func newPlayer(playerId uint64, clientAddr, gatewayGRPCAddr string) *Player {
	player := &Player{
		playerChanMsg:   make(chan []byte, 64),
		closeChannel:    make(chan bool, 0),
		gatewayGRPCAddr: gatewayGRPCAddr,
		clientAddr:      clientAddr,
		playerId:        playerId,
		innerRouter:     router.NewRouter(),
	}
	go player.loop()
	return player
}

func (player *Player) loop() {
	if !player.isInitFinish {
		player.init()
		player.isInitFinish = true
	}
	defer func() {
		WorldGetMe().delPlayer(player.playerId)
		logger.Infof("玩家(%v)退出online", player.playerId)
	}()
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

func (player *Player) addCloseChannel() {
	player.closeChannel <- true
}

func (player *Player) onMessage(data []byte) {
	msgID, err := player.innerRouter.Route(data)
	if err != nil {
		player.sendMsgToGateway(data)
		logger.Infof("其它服数据直接转发到gateway MsgId=%v", msgID)
	}
}

func (player *Player) init() {
	player.regMsgHandler()
	player.initBase()
}

func (player *Player) initBase() error {
	//userManger, err := NewUserManager(account, password)
	//if err != nil {
	//	return errors.New("")
	//}
	//player.User = userManger
	return nil
}

func (player *Player) sendMSg(cmd command.Command, msg interface{}) {
	data, err := player.innerRouter.Marshal(uint16(cmd), msg)
	if err != nil {
		return
	}
	player.sendMsgToGateway(data)
}

//数据转发到gateway
func (player *Player) sendMsgToGateway(data []byte) {
	gatewayStream := tcpServer().getGatewayStream(player.gatewayGRPCAddr)
	if gatewayStream == nil {
		return
	}
	gatewayStream.addToGatewayMsg(data, player.clientAddr)

	//
	player.sendMsgToPvp(data)
}

//数据转发到pvp
func (player *Player) sendMsgToPvp(data []byte) {
	pvpStream := PvpStreamGetMe()
	if pvpStream == nil {
		return
	}
	pvpStream.addPvpMsgChannel(player.playerId, data)
}
