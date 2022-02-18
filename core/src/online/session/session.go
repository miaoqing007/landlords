package session

import (
	"landlords/manager"
	"landlords/registry"
)

type Session struct {
	ch chan []byte
	*manager.Player
}

func NewSession(ch chan []byte) *Session {
	s := &Session{}
	s.ch = ch
	return s
}

//初始玩玩家信息
func (s *Session) InitPlayer(account, password string) {
	s.Player = &manager.Player{}
	if err := s.InitBase(account, password); err != nil {
		return
	}
	s.AddPlayer(s.Player.User.Id, s.Player)
	registry.Register(s.Player.User.Id, s.ch)
}

//玩家离线
func (s *Session) OffLine() {
	if s == nil || s.Player == nil || s.Player.User == nil {
		return
	}
	manager.RemoveRoom(s.User.GetRoomId())
	manager.RemovePlayer4PvpPool(s.User.GetPiecewise(), s.User.Id)
	manager.DeletePlayer(s.User.Id)
	registry.UnRegister(s.User.Id)
}
