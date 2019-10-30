package manager

import (
	"app/client_proto"
	"app/helper/common"
	"app/helper/conv"
	"app/operatecard"
	"sync"
)

var room *Rooms

type Rooms struct {
	rooms sync.Map //map[roomId]*RoomManager
}

func InitRoomManager() {
	room = &Rooms{}
	room.initCreateRoom()
}

func GetRoomManager(roomId string) *RoomManager {
	v, ok := room.rooms.Load(roomId)
	if ok {
		return v.(*RoomManager)
	}
	return nil
}

func (r *Rooms) initCreateRoom() {
	for i := 0; i < 10; i++ {
		r.initSingleRoom(conv.FormatInt(i))
	}
}

func (r *Rooms) initSingleRoom(roomId string) {
	r.rooms.Store(roomId, NewRoomManager(roomId))
}

type RoomManager struct {
	roomId               string
	holeCards            []string
	lastPlayerWasteCards sync.Map //map[uid][]cards
	player               sync.Map //map[uid]*UserInfo
}

func NewRoomManager(roomId string) *RoomManager {
	rm := &RoomManager{}
	rm.roomId = roomId
	return rm
}

func (r *RoomManager) ResetRoomManager() {
	r.holeCards = []string{}
	r.lastPlayerWasteCards = sync.Map{}
	r.player.Range(func(key, value interface{}) bool {
		value.(*UserInfo).resetUserInfo()
		return true
	})
}

func (r *RoomManager) GetUserInfo(uid string) *UserInfo {
	if u, ok := r.player.Load(uid); ok {
		return u.(*UserInfo)
	}
	return nil
}

func (r *RoomManager) CreatePlayerCards(cards1, cards2, cards3, holeCards []string, info *client_proto.S_player_cards) {
	num := 0
	r.player.Range(func(key, value interface{}) bool {
		if num%3 == 0 {
			value.(*UserInfo).createCards(cards1)
			info.F_player_1 = cards1
		} else if num%3 == 1 {
			value.(*UserInfo).createCards(cards2)
			info.F_player_2 = cards2
		} else if num%3 == 2 {
			value.(*UserInfo).createCards(cards3)
			info.F_plyaer_3 = cards3
		}
		num++
		return true
	})
	info.F_hole_cards = holeCards
}

func (r *RoomManager) CheckHandCards(uid string, cards []string) bool {
	v, ok := r.player.Load(uid)
	if !ok {
		return false
	}
	if !v.(*UserInfo).checkCards(cards) {
		return false
	}
	if !r.checkLastPlayerWasteCards(uid, cards) {
		return false
	}
	return true
}

func (r *RoomManager) checkLastPlayerWasteCards(uid string, cards []string) bool {
	lastWasteCards, ok := r.lastPlayerWasteCards.Load(uid)
	if ok || common.GetSyncMapLen(r.lastPlayerWasteCards) == 0 {
		r.updateLastPlayerWasteCards(uid, cards)
		return true
	} else if operatecard.ComparisonTwoPlayersCards(lastWasteCards.([]string), cards) {
		r.lastPlayerWasteCards = sync.Map{}
		r.updateLastPlayerWasteCards(uid, cards)
		return true
	} else {
		return false
	}
}

func (r *RoomManager) updateLastPlayerWasteCards(uid string, cards []string) {
	r.lastPlayerWasteCards.Store(uid, cards)
}

func (r *RoomManager) UpdateCards(uid string, useCards []string) bool {
	v, ok := r.player.Load(uid)
	if !ok {
		return false
	}
	v.(*UserInfo).updateCards(useCards)
	return true
}

func (r *RoomManager) AddPlayerToRoom(uid string) bool {
	if common.GetSyncMapLen(r.player) >= 3 {
		return false
	}
	r.player.Store(uid, NewUserInfo(uid))
	return true
}

type UserInfo struct {
	id        string
	ordinal   int
	handCards sync.Map //map[card]bool
}

func NewUserInfo(uid string) *UserInfo {
	u := &UserInfo{}
	u.id = uid
	return u
}

func (u *UserInfo) updateCards(wasteCards []string) {
	for _, card := range wasteCards {
		u.handCards.Delete(card)
	}
}

func (u *UserInfo) createCards(cards []string) {
	for _, v := range cards {
		u.handCards.Store(v, true)
	}
}

func (u *UserInfo) checkCards(cards []string) bool {
	for _, card := range cards {
		if _, ok := u.handCards.Load(card); !ok {
			return ok
		}
	}
	return true
}

func (u *UserInfo) resetUserInfo() {
	u.handCards = sync.Map{}
	u.ordinal = 0
}
