package manager

import (
	"github.com/golang/glog"
	"landlords/client_proto"
	"landlords/enmu"
	"landlords/helper/common"
	"landlords/helper/conv"
	"landlords/helper/util"
	"landlords/initcards"
	. "landlords/obj"
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
func Add2Room(piecewise int, pms []AddChanMsg) {
	uids := make([]string, 0)
	for _, pm := range pms {
		uids = append(uids,pm.Id )
	}
	rm := NewRoomManager(piecewise, pms)
	registry.RegisterRoom(rm.roomId, uids)
	room.rooms.Store(rm.roomId, rm)
	push2Client(rm.roomId, uids)
	glog.Infof("roomId = %v ,uids = %v", rm.roomId, pms)
}

func push2Client(roomId string, ids []string) {
	info := client_proto.S_player_card{}
	cards := initcards.ShuffCards()
	room := GetRoomManager(roomId)
	if room == nil {
		return
	}
	room.CreatePlayerCards(cards, &info, ids)
}

type RoomManager struct {
	roomId               string   //房间id
	holeCards            []string //底牌
	piecewise            int      //分段
	landowner            string   //地主id
	lastPlayerWasteCards sync.Map //map[uid][]cards
	players              sync.Map //map[uid]*UserInfo
}

func NewRoomManager(piecewise int, pms []AddChanMsg) *RoomManager {
	rm := &RoomManager{}
	rm.piecewise = piecewise
	rm.roomId = conv.FormatInt32(atomic.AddInt32(&room.roomIcreseId, 1))
	for _, pm := range pms {
		rm.AddPlayerToRoom(pm.Id, pm.Name)
	}
	return rm
}

//重置房间信息
func (r *RoomManager) ResetRoomManager() {
	r.holeCards = []string{}
	r.lastPlayerWasteCards = sync.Map{}
	r.players.Range(func(key, value interface{}) bool {
		value.(*UserInfo).resetUserInfo()
		return true
	})
}

func (r *RoomManager) RemoveManager() {
	r.players.Range(func(key, value interface{}) bool {
		r.players.Delete(key)
		return true
	})
	registry.PushRoom(r.roomId, 2001, client_proto.S_entity_id{})
	registry.UnRegisterRoom(r.roomId)
}

//获取玩家信息
func (r *RoomManager) GetUserInfo(uid string) *UserInfo {
	if u, ok := r.players.Load(uid); ok {
		return u.(*UserInfo)
	}
	return nil
}

//创建玩家手牌
func (r *RoomManager) CreatePlayerCards(cards []string, info *client_proto.S_player_card, ids []string) {
	r.players.Range(func(key, value interface{}) bool {
		value.(*UserInfo).addCards(cards[:17])
		info.F_players = append(info.F_players, client_proto.S_player{value.(*UserInfo).id,
			value.(*UserInfo).name, util.SortArrayString(cards[:17])})
		cards = cards[17:]
		return true
	})
	info.F_hole_cards = util.SortArrayString(cards)
	info.F_playerIds = ids
	info.F_roomId = r.roomId
	registry.PushRoom(info.F_roomId, 3000, info)
}

//判断玩家手牌
func (r *RoomManager) CheckHandCards(uid string, cards []string) bool {
	v, ok := r.players.Load(uid)
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
	cardType := operatecard.JudgeCardsType(operatecard.GetCardsValue(cards))
	if cardType == enmu.ERROR_TYPE {
		return false
	}
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
	v, ok := r.players.Load(uid)
	if !ok {
		return false
	}
	v.(*UserInfo).delCards(useCards)
	return true
}

//获取玩家剩余手牌数量
func (r *RoomManager) GetUserRemainingCardsNum(uid string) int {
	v, ok := r.players.Load(uid)
	if !ok {
		return -1
	}
	return v.(*UserInfo).getUserRemainingCardsNum()
}

func (r *RoomManager) GetUserCard(uid string) sync.Map {
	v, ok := r.players.Load(uid)
	if !ok {
		return sync.Map{}
	}
	return v.(*UserInfo).getUserCards()
}

func (r *RoomManager) GetUserCard4Array(uid string) []string {
	cards := make([]string, 0)
	cm := r.GetUserCard(uid)
	cm.Range(func(key, value interface{}) bool {
		cards = append(cards, key.(string))
		return true
	})
	return cards
}

//添加玩家到房间
func (r *RoomManager) AddPlayerToRoom(uid, name string) bool {
	if common.GetSyncMapLen(r.players) >= 3 {
		return false
	}
	p := GetPlayer(uid)
	if p != nil {
		p.User.SetRoomId(r.roomId)
	}
	r.players.Store(uid, NewUserInfo(uid, name))
	return true
}

type UserInfo struct {
	id                string //玩家id
	name              string
	handCards         sync.Map //map[card]bool 手牌
	remainingCardsNum int      //剩余手牌数量
}

func NewUserInfo(uid, name string) *UserInfo {
	u := &UserInfo{}
	u.id = uid
	u.name = name
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
