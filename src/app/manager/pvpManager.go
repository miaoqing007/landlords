package manager

import (
	"github.com/golang/glog"
	"sync"
	"time"
)

var _pvpPoolManger *PvpPoolManager

type PvpPoolManager struct {
	ups     sync.Map //map[int分段][]string玩家id
	addChan chan addChanMsg
	remChan chan remChanMsg
}

type addChanMsg struct {
	piecewise int
	id        string
}

type remChanMsg struct {
	piecewise int
	id        string
}

//初始化匹配池
func InitPvpPoolManager() {
	_pvpPoolManger = &PvpPoolManager{}
	_pvpPoolManger.addChan = make(chan addChanMsg, 64)
	_pvpPoolManger.remChan = make(chan remChanMsg, 64)
	go _pvpPoolManger.watch()
	glog.Info("初始pvp完成")
}

//监听匹配相关信息
func (p *PvpPoolManager) watch() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case acm := <-p.addChan:
			p.addups(acm)
		case dcm := <-p.remChan:
			p.remups(dcm)
		case <-ticker.C:
			p.pvpMatchPlayer()
		default:
		}
	}
}

//匹配玩家
func (p *PvpPoolManager) pvpMatchPlayer() {
	p.ups.Range(func(key, value interface{}) bool {
		arr := value.([]string)
		for len(arr) >= 3 {
			us := arr[:3]
			arr = arr[3:]
			Add2Room(key.(int), us)
		}
		p.ups.Store(key, arr)
		return true
	})
}

//添加玩家
func (p *PvpPoolManager) addups(acm addChanMsg) {
	ps, ok := p.ups.Load(acm.piecewise)
	if !ok {
		ps = make([]string, 0)
	}
	arr := ps.([]string)
	arr = append(arr, acm.id)
	p.ups.Store(acm.piecewise, arr)
}

//移除玩家
func (p *PvpPoolManager) remups(rcm remChanMsg) {
	ps, ok := p.ups.Load(rcm.piecewise)
	if !ok {
		return
	}
	arr := ps.([]string)
	for k, id := range arr {
		if id == rcm.id {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	p.ups.Store(rcm.piecewise, arr)
}

//玩家进入匹配池
func AddPlayer2PvpPool(piecewise int, id string) {
	_pvpPoolManger.addChan <- addChanMsg{piecewise, id}
}

//移除匹配池的玩家
func RemovePlayer4PvpPool(piecewise int, id string) {
	_pvpPoolManger.remChan <- remChanMsg{piecewise, id}
}
