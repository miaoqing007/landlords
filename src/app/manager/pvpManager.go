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
	delChan chan delChanMsg
}

type addChanMsg struct {
	piecewise int
	id        string
}

type delChanMsg struct {
	piecewise int
	id        string
}

func InitPvpPoolManager() {
	_pvpPoolManger = &PvpPoolManager{}
	_pvpPoolManger.addChan = make(chan addChanMsg, 64)
	_pvpPoolManger.delChan = make(chan delChanMsg, 64)
	go _pvpPoolManger.watch()
	glog.Info("初始pvp完成")
}

func (p *PvpPoolManager) watch() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case acm := <-p.addChan:
			p.addups(acm)
		case dcm := <-p.delChan:
			p.delups(dcm)
		case <-ticker.C:
			p.pvpMatchPlayer()
		default:
		}
	}
}

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

func (p *PvpPoolManager) addups(acm addChanMsg) {
	ps, ok := p.ups.Load(acm.piecewise)
	if !ok {
		ps = make([]string, 0)
	}
	arr := ps.([]string)
	arr = append(arr, acm.id)
	p.ups.Store(acm.piecewise, arr)
}

func (p *PvpPoolManager) delups(dcm delChanMsg) {
	ps, ok := p.ups.Load(dcm.piecewise)
	if !ok {
		return
	}
	arr := ps.([]string)
	for k, id := range arr {
		if id == dcm.id {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	p.ups.Store(dcm.piecewise, arr)
}

func AddPlayer2PvpPool(piecewise int, id string) {
	_pvpPoolManger.addChan <- addChanMsg{piecewise, id}
}

func DelPlayer4PvpPool(piecewise int, id string) {
	_pvpPoolManger.delChan <- delChanMsg{piecewise, id}
}
