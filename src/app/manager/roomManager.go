package manager

import (
	"app/client_proto"
	"app/helper/common"
	"app/helper/conv"
	"app/misc/packet"
	"app/operatecard"
	"app/registry"
	"github.com/golang/glog"
	"sync"
	"sync/atomic"
)

var room *Rooms

type Rooms struct {
	roomIcreseId int32    //递增房间id
	rooms        sync.Map //map[roomId]*RoomManager
}

func InitRoomManager() {
	room = &Rooms{}
	glog.Info("初始化房间完成")
}

//获取房间信息
func GetRoomManager(roomId string) *RoomManager {
	v, ok := room.rooms.Load(roomId)
	if ok {
		return v.(*RoomManager)
	}
	return nil
}

func RemoveRoom(roomId string) {
	registry.UnRegisterRoom(roomId)
	room.rooms.Delete(roomId)
}

//添加玩家到房间
func Add2Room(piecewise int, ids []string) {
	rm := NewRoomManager(piecewise, ids)
	registry.RegisterRoom(rm.roomId, ids)
	room.rooms.Store(rm.roomId, rm)
	glog.Infof("roomId = %v ,uids = %v", rm.roomId, ids)
}

type RoomManager struct {
	roomId               string   //房间id
	holeCards            []string //底牌
	piecewise            int      //分段
	lastPlayerWasteCards sync.Map //map[uid][]cards
	player               sync.Map //map[uid]*UserInfo
}

func NewRoomManager(piecewise int, ids []string) *RoomManager {
	rm := &RoomManager{}
	rm.piecewise = piecewise
	rm.roomId = conv.FormatInt32(atomic.AddInt32(&room.roomIcreseId, 1))
	for _, id := range ids {
		rm.AddPlayerToRoom(id)
	}
	return rm
}

//重置房间信息
func (r *RoomManager) ResetRoomManager() {
	r.holeCards = []string{}
	r.lastPlayerWasteCards = sync.Map{}
	r.player.Range(func(key, value interface{}) bool {
		value.(*UserInfo).resetUserInfo()
		return true
	})
}

//获取玩家信息
func (r *RoomManager) GetUserInfo(uid string) *UserInfo {
	if u, ok := r.player.Load(uid); ok {
		return u.(*UserInfo)
	}
	return nil
}

//创建玩家手牌
func (r *RoomManager) CreatePlayerCards(cards1, cards2, cards3, holeCards []string, info *client_proto.S_player_card) {
	num := 0
	r.player.Range(func(key, value interface{}) bool {
		if num%3 == 0 {
			value.(*UserInfo).addCards(cards1)
			info.F_players = append(info.F_players, client_proto.S_player{value.(*UserInfo).id, cards1})
		} else if num%3 == 1 {
			value.(*UserInfo).addCards(cards2)
			info.F_players = append(info.F_players, client_proto.S_player{value.(*UserInfo).id, cards2})
		} else if num%3 == 2 {
			value.(*UserInfo).addCards(cards3)
			info.F_players = append(info.F_players, client_proto.S_player{value.(*UserInfo).id, cards3})
		}
		num++
		return true
	})
	info.F_hole_cards = holeCards
	info.F_roomId = r.roomId
	registry.PushRoom(info.F_roomId, packet.Pack(2003, info, nil))
}

//判断玩家手牌
func (r *RoomManager) CheckHandCards(uid string, cards []string) bool {
	v, ok := r.player.Load(uid)
	if !ok {
		return false
	}
	if !v.(*UserInfo).checkCards(cards) {
		return false
	}
	if !r.comparisonLastPlayerWasteCards(uid, cards) {
		return false
	}
	return true
}

//比对玩家手牌
func (r *RoomManager) comparisonLastPlayerWasteCards(uid string, cards []string) bool {
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

//更新上一家所出的牌
func (r *RoomManager) updateLastPlayerWasteCards(uid string, cards []string) {
	r.lastPlayerWasteCards.Store(uid, cards)
}

//删除玩家已出手牌
func (r *RoomManager) DeleteCards(uid string, useCards []string) bool {
	v, ok := r.player.Load(uid)
	if !ok {
		return false
	}
	v.(*UserInfo).delCards(useCards)
	return true
}

//获取玩家剩余手牌数量
func (r *RoomManager) GetUserRemainingCardsNum(uid string) int {
	v, ok := r.player.Load(uid)
	if !ok {
		return -1
	}
	return v.(*UserInfo).getUserRemainingCardsNum()
}

//添加玩家到房间
func (r *RoomManager) AddPlayerToRoom(uid string) bool {
	if common.GetSyncMapLen(r.player) >= 3 {
		return false
	}
	r.player.Store(uid, NewUserInfo(uid))
	return true
}

type UserInfo struct {
	id                string   //玩家id
	handCards         sync.Map //map[card]bool 手牌
	remainingCardsNum int      //剩余手牌数量
}

func NewUserInfo(uid string) *UserInfo {
	u := &UserInfo{}
	u.id = uid
	return u
}

//删除已出手牌
func (u *UserInfo) delCards(wasteCards []string) {
	for _, card := range wasteCards {
		u.handCards.Delete(card)
	}
	u.remainingCardsNum = common.GetSyncMapLen(u.handCards)
}

//增加手牌
func (u *UserInfo) addCards(cards []string) {
	for _, v := range cards {
		u.handCards.Store(v, true)
	}
	u.remainingCardsNum = common.GetSyncMapLen(u.handCards)
}

//核对手牌
func (u *UserInfo) checkCards(cards []string) bool {
	for _, card := range cards {
		if _, ok := u.handCards.Load(card); !ok {
			return ok
		}
	}
	return true
}

//获取剩余手牌数量
func (u *UserInfo) getUserRemainingCardsNum() int {
	return u.remainingCardsNum
}

//重置信息
func (u *UserInfo) resetUserInfo() {
	u.handCards = sync.Map{}
}
