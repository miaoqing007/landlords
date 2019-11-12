package session

import (
	"app/manager"
	"app/misc/packet"
	"app/registry"
	"github.com/golang/glog"
)

var ExecuteHandler func(code int16, sess *Session, reader *packet.Packet) [][]byte

type Session struct {
	die         chan struct{}
	recieveChan chan []byte
	sendChan    chan []byte
	rSendChan   chan []byte
	*manager.Player
}

func NewSession() *Session {
	s := &Session{}
	s.die = make(chan struct{}, 1)
	s.recieveChan = make(chan []byte, 16)
	s.sendChan = make(chan []byte, 16)
	s.rSendChan = make(chan []byte, 16)
	go s.watch()
	return s
}

func (s *Session) watch() {
	for {
		select {
		case <-s.die:
			s.OffLine(s.Id)
			return
		case msg := <-s.recieveChan:
			reader := packet.Reader(msg)
			c, err := reader.ReadS16()
			if err != nil {
				glog.Info("err=", err)
				return
			}
			bytes := ExecuteHandler(c, s, reader)
			for _, byt := range bytes {
				s.sendChan <- byt
			}
		case msg := <-s.rSendChan:
			s.sendChan <- msg
		default:
		}
	}
}

//注册接收信息channel
func (s *Session) EvaluationReciveChan(ch chan []byte) {
	s.recieveChan = ch
}

//注册发送信息channel
func (s *Session) EvaluationSendChan(ch chan []byte) {
	s.sendChan = ch
}

func (s *Session) AddDieChan() {
	s.die <- struct{}{}
}

//初始玩玩家信息
func (s *Session) InitPlayer(id string) error {
	s.Player = &manager.Player{}
	userManger, err := manager.NewUserManager(id)
	if err != nil {
		return err
	}
	s.UserManager = userManger

	manager.AddPlayer(s.Player.UserManager.Id, s.Player)
	registry.Register(id, s.rSendChan)
	return nil
}

//玩家离线
func (s *Session) OffLine(id string) {
	manager.DeletePlayer(id)
	registry.UnRegister(id)
}
