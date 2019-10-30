package registry

import (
	"app/session"
	"sync"
)

var onlineUser *Registry

type Registry struct {
	users sync.Map
}

func init() {
	onlineUser = &Registry{}
}

func (r *Registry) pushMSg(uid string, msg []byte) {
	v, ok := r.users.Load(uid)
	if !ok {
		return
	}
	v.(*session.Session).AddSendChan(msg)
}

func (r *Registry) registry(uid string, sess *session.Session) {
	r.users.Store(uid, sess)
}

func (r *Registry) unRegistry(uid string) {
	r.users.Delete(uid)
}

func Register(uid string, sess *session.Session) {
	onlineUser.registry(uid, sess)
}

func Push(uid string, msg []byte) {
	onlineUser.pushMSg(uid, msg)
}

func UnRegister(uid string) {
	onlineUser.unRegistry(uid)
}
