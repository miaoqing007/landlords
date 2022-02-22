package main

//
//import (
//	"github.com/golang/glog"
//	"sync"
//	"time"
//)
//
//var _pvpPoolManger *PvpPoolManager
//
//type PvpPoolManager struct {
//	ups     sync.Map        //map[int分段][]addChanMsg玩家id
//	addChan chan AddChanMsg //添加匹配池channel
//	remChan chan remChanMsg //移除匹配池channel
//}
//
////type addChanMsg struct {
////	piecewise int
////	id        string
////	name      string
////}
//
//type remChanMsg struct {
//	piecewise int
//	id        string
//}
//
////初始化匹配池
//func InitPvpPoolManager() {
//	_pvpPoolManger = &PvpPoolManager{}
//	_pvpPoolManger.addChan = make(chan AddChanMsg, 64)
//	_pvpPoolManger.remChan = make(chan remChanMsg, 64)
//	go _pvpPoolManger.watch()
//	glog.Info("初始pvp完成")
//}
//
////监听匹配相关信息
//func (p *PvpPoolManager) watch() {
//	ticker := time.NewTicker(1 * time.Second)
//	defer ticker.Stop()
//	for {
//		select {
//		case acm := <-p.addChan:
//			p.addups(acm)
//		case dcm := <-p.remChan:
//			p.remups(dcm)
//		case <-ticker.C:
//			p.pvpMatchPlayer()
//		}
//	}
//}
//
////匹配玩家
//func (p *PvpPoolManager) pvpMatchPlayer() {
//	p.ups.Range(func(key, value interface{}) bool {
//		arr := value.([]AddChanMsg)
//		if len(arr) < 3 {
//			return false
//		}
//		for len(arr) >= 3 {
//			us := arr[:3]
//			arr = arr[3:]
//			Add2Room(key.(int), us)
//		}
//		p.ups.Store(key, arr)
//		return true
//	})
//}
//
////添加玩家
//func (p *PvpPoolManager) addups(acm AddChanMsg) {
//	ps, ok := p.ups.Load(acm.Piecewise)
//	if !ok {
//		ps = make([]AddChanMsg, 0)
//	}
//	arr := ps.([]AddChanMsg)
//	for _, v := range arr {
//		if v.Id == acm.Id {
//			return
//		}
//	}
//	player := GetPlayer(acm.Id)
//	if player != nil {
//		player.User.SetPiecewise(acm.Piecewise)
//	}
//	arr = append(arr, acm)
//	p.ups.Store(acm.Piecewise, arr)
//}
//
////移除玩家
//func (p *PvpPoolManager) remups(rcm remChanMsg) {
//	ps, ok := p.ups.Load(rcm.piecewise)
//	if !ok {
//		return
//	}
//	arr := ps.([]AddChanMsg)
//	for k, pm := range arr {
//		if pm.Id == rcm.id {
//			arr = append(arr[:k], arr[k+1:]...)
//			break
//		}
//	}
//	p.ups.Store(rcm.piecewise, arr)
//}
//
////玩家进入匹配池
//func (player *Player)addPlayer2PvpPool(piecewise int, id, name string) {
//	_pvpPoolManger.addChan <- AddChanMsg{id, name, piecewise}
//	glog.Infof("addPlayerPvpPool = %v", id)
//}
//
////移除匹配池的玩家
//func RemovePlayer4PvpPool(piecewise int, id string) {
//	_pvpPoolManger.remChan <- remChanMsg{piecewise, id}
//	glog.Infof("removePlayerPvpPool = %v", id)
//}
