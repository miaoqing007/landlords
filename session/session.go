package session

import (
	"landlords/manager"
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

////初始玩玩家信息
//func (s *Session) InitPlayer(id string) error {
//	s.Player = &manager.Player{}
//	userManger, err := manager.NewUserManager(id)
//	if err != nil {
//		return err
//	}
//	s.User = userManger
//
//	manager.AddPlayer(s.User.Account, s.Player)
//	registry.Register(s.User.Account, s.ch)
//	return nil
//}
//
////玩家离线
//func (s *Session) OffLine(id string) {
//	manager.RemoveRoom(s.User.GetRoomId())
//	manager.RemovePlayer4PvpPool(s.User.GetPiecewise(), id)
//	manager.DeletePlayer(id)
//	registry.UnRegister(id)
//}
