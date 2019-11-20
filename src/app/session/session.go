package session

import (
	"app/manager"
	"app/registry"
)

type Session struct {
	Id string
	ch chan []byte
	*manager.Player
}

func NewSession(ch chan []byte) *Session {
	s := &Session{}
	s.ch = ch
	return s
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
	registry.Register(id, s.ch)
	return nil
}

//玩家离线
func (s *Session) OffLine(id string) {
	manager.DeletePlayer(id)
	registry.UnRegister(id)
}
