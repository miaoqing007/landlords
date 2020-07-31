package registry

import (
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"landlords/misc/packet"
	. "landlords/obj"
	"sync"
)

var (
	onlineUser *Registry //在线玩家注册信息
)

type Registry struct {
	users sync.Map      //map[id]regmsg
	rooms sync.Map      //map[roomid][]string
	rch   chan regmsg   //接收注册信息channel
	urch  chan string   //接收反注册channel
	pch   chan pushmsg  //接收消息推送channel
	rRch  chan roommsg  //接收房间注册channel
	urRch chan string   //接收房间反注册channel
	rpch  chan rpushmsg //接收房间消息推送channel
}

type pushmsg struct {
	uid string
	msg *WsMessage
}

type regmsg struct {
	uid    string
	sendch chan *WsMessage
}

type roommsg struct {
	roomid string
	uids   []string
}

type rpushmsg struct {
	roomid string
	msg    *WsMessage
}

func init() {
	onlineUser = &Registry{}
	onlineUser.rch = make(chan regmsg, 16)
	onlineUser.urch = make(chan string, 16)
	onlineUser.pch = make(chan pushmsg, 16)
	onlineUser.rRch = make(chan roommsg, 16)
	onlineUser.urRch = make(chan string, 16)
	onlineUser.rpch = make(chan rpushmsg, 16)
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
		case rRmsg := <-r.rRch:
			r.registryRoom(rRmsg)
		case rid := <-r.urRch:
			r.unRegistryRoom(rid)
		case rpmsg := <-r.rpch:
			r.rpushMsg(rpmsg)
		}
	}
}

func (r *Registry) pushMSg(pmsg pushmsg) {
	v, ok := r.users.Load(pmsg.uid)
	if !ok {
		return
	}
	v.(regmsg).sendch <- pmsg.msg
}

func (r *Registry) rpushMsg(rpmsg rpushmsg) {
	v, ok := r.rooms.Load(rpmsg.roomid)
	if !ok {
		return
	}
	for _, id := range v.([]string) {
		Push(id, rpmsg.msg)
	}
}

func (r *Registry) registry(rm regmsg) {
	r.users.Store(rm.uid, rm)
}

func (r *Registry) registryRoom(rR roommsg) {
	r.rooms.Store(rR.roomid, rR.uids)
}

func (r *Registry) unRegistry(uid string) {
	r.users.Delete(uid)
}

func (r *Registry) unRegistryRoom(roomId string) {
	r.rooms.Delete(roomId)
}

//玩家注册
func Register(uid string, sch chan *WsMessage) {
	onlineUser.rch <- regmsg{uid: uid, sendch: sch}
	glog.Infof("register = %v", uid)
}

//玩家消息推送
func Push(uid string, msg *WsMessage) {
	onlineUser.pch <- pushmsg{uid, msg}
}

//玩家反注册
func UnRegister(uid string) {
	onlineUser.urch <- uid
	glog.Infof("unregister = %v", uid)
}

//房间注册
func RegisterRoom(roomid string, uids []string) {
	onlineUser.rRch <- roommsg{roomid, uids}
	glog.Infof("registerroom = %v uids = %v", roomid, uids)
}

//房间反注册
func UnRegisterRoom(roomid string) {
	onlineUser.urRch <- roomid
	glog.Infof("unregisterroom = %v", roomid)
}

//房间消息推送
func PushRoom(roomid string, tos int16, ret interface{}) {
	onlineUser.rpch <- rpushmsg{roomid: roomid,
		msg: &WsMessage{MessageType: websocket.TextMessage,
			Data: packet.Pack(tos, ret)}}
}
