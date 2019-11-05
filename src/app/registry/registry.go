package registry

import (
	"sync"
)

var onlineUser *Registry

type Registry struct {
	users sync.Map //map[id]regMsg
	rch   chan regMsg
	urch  chan string
	pch   chan pushmsg
}

type pushmsg struct {
	uid string
	msg []byte
}

type regMsg struct {
	uid    string
	sendch chan []byte
}

func init() {
	onlineUser = &Registry{}
	onlineUser.rch = make(chan regMsg, 16)
	onlineUser.urch = make(chan string, 16)
	onlineUser.pch = make(chan pushmsg, 16)
	go onlineUser.watch()
}

func (r *Registry) watch() {
	for {
		select {
		case rm := <-r.rch:
			r.registry(rm)
		case id := <-r.urch:
			r.unRegistry(id)
		case pmsg := <-r.pch:
			r.pushMSg(pmsg)
		default:
		}
	}
}

func (r *Registry) pushMSg(pmsg pushmsg) {
	v, ok := r.users.Load(pmsg.uid)
	if !ok {
		return
	}
	v.(regMsg).sendch <- (pmsg.msg)
}

func (r *Registry) registry(rm regMsg) {
	r.users.Store(rm.uid, rm)
}

func (r *Registry) unRegistry(uid string) {
	r.users.Delete(uid)
}

func Register(uid string, sch chan []byte) {
	onlineUser.rch <- regMsg{uid, sch}
}

func Push(uid string, msg []byte) {
	onlineUser.pch <- pushmsg{uid, msg}
}

func UnRegister(uid string) {
	onlineUser.urch <- uid
}
