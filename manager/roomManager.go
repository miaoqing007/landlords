package manager

import (
	"github.com/golang/glog"
	"landlords/client_proto"
	"landlords/helper/common"
	"landlords/helper/conv"
	"landlords/initcards"
	"landlords/operatecard"
	"landlords/registry"
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

//移除房间
func RemoveRoom(roomId string) {
	r, ok := room.rooms.Load(roomId)
	if !ok {
		return
	}
	r.(*RoomManager).RemoveManager()
	room.rooms.Delete(roomId)
	glog.Infof("remove roomId = %v", roomId)
}

//添加玩家到房间
func Add2Room(piecewise int, ids []string) {
	rm := NewRoomManager(piecewise, ids)
	registry.RegisterRoom(rm.roomId, ids)
	room.rooms.Store(rm.roomId, rm)
	T(rm.roomId)
	glog.Infof("roomId = %v ,uids = %v", rm.roomId, ids)
}

func T(roomId string) {
	info := client_proto.S_player_card{}
	cards := initcards.ShuffCards()
	room := GetRoomManager(roomId)
	if room == nil {
		return
	}
	room.CreatePlayerCards(cards, &info)
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

func (r *RoomManager) RemoveManager() {
	r.player.Range(func(key, value interface{}) bool {
		r.player.Delete(key)
		return true
	})
	registry.PushRoom(r.roomId, 2001, client_proto.S_entity_id{})
	registry.UnRegisterRoom(r.roomId)
}

//获取玩家信息
func (r *RoomManager) GetUserInfo(uid string) *UserInfo {
	if u, ok := r.player.Load(uid); ok {
		return u.(*UserInfo)
	}
	return nil
}

//创建玩家手牌
func (r *RoomManager) CreatePlayerCards(cards []string, info *client_proto.S_player_card) {
	r.player.Range(func(key, value interface{}) bool {
		value.(*UserInfo).addCards(cards[:17])
		info.F_players = append(info.F_players, client_proto.S_player{value.(*UserInfo).id, cards[:17]})
		cards = cards[17:]
		return true
	})
	info.F_hole_cards = cards
	info.F_roomId = r.roomId
	registry.PushRoom(info.F_roomId, 2003, info)
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
	_, ok := r.lastPlayerWasteCards.Load(uid)
	if ok || common.GetSyncMapLen(r.lastPlayerWasteCards) == 0 {
		r.updateLastPlayerWasteCards(uid, cards)
		return true
	} else {
		var lastWasteCards []string
		r.lastPlayerWasteCards.Range(func(key, value interface{}) bool {
			lastWasteCards = value.([]string)
			return true
		})
		if operatecard.ComparisonTwoPlayersCards(lastWasteCards, cards) {
			r.lastPlayerWasteCards = sync.Map{}
			r.updateLastPlayerWasteCards(uid, cards)
			return true
		} else {
			return false
		}
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

func (r *RoomManager) GetUserCard(uid string) sync.Map {
	v, ok := r.player.Load(uid)
	if !ok {
		return sync.Map{}
	}
	return v.(*UserInfo).getUserCards()
}

//添加玩家到房间
func (r *RoomManager) AddPlayerToRoom(uid string) bool {
	if common.GetSyncMapLen(r.player) >= 3 {
		return false
	}
	p := GetPlayer(uid)
	if p != nil {
		p.User.SetRoomId(r.roomId)
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

func (u *UserInfo) getUserCards() sync.Map {
	return u.handCards
}

//重置信息
func (u *UserInfo) resetUserInfo() {
	u.handCards = sync.Map{}
}
