package registry

import (
	"sync"
)

var onlineUser *Registry //在线玩家注册信息

type Registry struct {
	users sync.Map     //map[id]regMsg
	rch   chan regMsg  //接收注册信息channel
	urch  chan string  //接收反注册channel
	pch   chan pushmsg //接收消息推送channel
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
	v.(regMsg).sendch <- pmsg.msg
}

func (r *Registry) registry(rm regMsg) {
	r.users.Store(rm.uid, rm)
}

func (r *Registry) unRegistry(uid string) {
	r.users.Delete(uid)
}

//玩家注册
func Register(uid string, sch chan []byte) {
	onlineUser.rch <- regMsg{uid, sch}
}

//消息推送
func Push(uid string, msg []byte) {
	onlineUser.pch <- pushmsg{uid, msg}
}

//玩家反注册
func UnRegister(uid string) {
	onlineUser.urch <- uid
}
