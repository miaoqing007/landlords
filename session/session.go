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
func (s *Session) InitPlayer(account, password string) error {
	s.Player = &manager.Player{}
	userManger, err := manager.NewUserManager(account, password)
	if err != nil {
		return err
	}
	s.User = userManger

	manager.AddPlayer(s.User.Id, s.Player)
	registry.Register(s.User.Id, s.ch)
	return nil
}

//玩家离线
func (s *Session) OffLine() {
	if s == nil || s.Player == nil {
		return
	}
	manager.RemoveRoom(s.User.GetRoomId())
	manager.RemovePlayer4PvpPool(s.User.GetPiecewise(), s.User.Id)
	manager.DeletePlayer(s.User.Id)
	registry.UnRegister(s.User.Id)
}
