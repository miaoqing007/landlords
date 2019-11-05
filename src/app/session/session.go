package session

import (
	"app/manager"
	"app/misc/packet"
	"github.com/golang/glog"
)

var ExecuteHandler func(code int16, sess *Session, reader *packet.Packet) [][]byte

type Session struct {
	die         chan struct{}
	recieveChan chan []byte
	sendChan    chan []byte
	*manager.Player
}

func NewSession() *Session {
	s := &Session{}
	s.die = make(chan struct{}, 1)
	s.recieveChan = make(chan []byte, 1)
	s.sendChan = make(chan []byte, 1)
	go s.watch()
	return s
}

func (s *Session) watch() {
	for {
		select {
		case <-s.die:
			s.OffLine(s.Id)
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
		default:
		}
	}
}

func (s *Session) AddRecieveChan(byte []byte) {
	s.recieveChan <- byte
}

func (s *Session) AddSendChan(byte []byte) {
	s.sendChan <- byte
}

func (s *Session) EvaluationReciveChan(ch chan []byte) {
	s.recieveChan = ch
}
func (s *Session) EvaluationSendChan(ch chan []byte) {
	s.sendChan = ch
}

func (s *Session) AddDieChan() {
	s.die <- struct{}{}
}

func (s *Session) InitUser(id string) error {
	s.Player = &manager.Player{}
	userManger, err := manager.NewUserManager(id)
	if err != nil {
		return err
	}
	s.UserManager = userManger

	manager.AddPlayer(s.Player.UserManager.Id, s.Player)
	return nil
}

func (s *Session) OffLine(id string) {
	manager.DeletePlayer(id)
}
