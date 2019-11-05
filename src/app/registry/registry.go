package registry

import (
	"app/session"
	"sync"
)

var onlineUser *Registry

type Registry struct {
	users sync.Map //map[id]*session
	rch   chan *session.Session
	urch  chan string
	pch   chan pushmsg
}

type pushmsg struct {
	uid string
	msg []byte
}

func init() {
	onlineUser = &Registry{}
	onlineUser.rch = make(chan *session.Session)
	onlineUser.urch = make(chan string)
	onlineUser.pch = make(chan pushmsg)
	go onlineUser.watch()
}

func (r *Registry) watch() {
	for {
		select {
		case sess := <-r.rch:
			r.registry(sess)
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
	v.(*session.Session).AddSendChan(pmsg.msg)
}

func (r *Registry) registry(sess *session.Session) {
	r.users.Store(sess.Id, sess)
}

func (r *Registry) unRegistry(uid string) {
	r.users.Delete(uid)
}

func Register(sess *session.Session) {
	onlineUser.rch <- sess
}

func Push(uid string, msg []byte) {
	onlineUser.pch <- pushmsg{uid, msg}
}

func UnRegister(uid string) {
	onlineUser.urch <- uid
}
