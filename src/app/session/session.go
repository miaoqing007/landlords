package session

import "app/manager"

type Session struct {
	Die         chan struct{}
	recieveChan chan []byte
	sendChan    chan []byte
	Player      *manager.Player
}

func NewSession() *Session {
	s := &Session{}
	s.Die = make(chan struct{}, 1)
	s.recieveChan = make(chan []byte, 1)
	s.sendChan = make(chan []byte, 1)
	return s
}

func (s *Session) AddRecieveChan(byte []byte) {
	s.recieveChan <- byte
}

func (s *Session) AddSendChan(byte []byte) {
	s.sendChan <- byte
}

func (s *Session) EvaluationSendChan(ch chan []byte) {
	s.sendChan = ch
}

func (s *Session) AddDieChan() {
	s.Die <- struct{}{}
}

func (s *Session) InitUser() {

}
