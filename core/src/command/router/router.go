package router

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"math"
)

type MsgHandler func(msgId uint16, data []byte)

type Router struct {
	littleEndian bool
	msgHandlers  map[uint16]MsgHandler
}

func NewRouter() *Router {
	r := &Router{
		littleEndian: true,
		msgHandlers:  make(map[uint16]MsgHandler),
	}
	return r
}

func (r *Router) Route(data []byte) (uint16, error) {
	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}
	var msgId uint16
	if r.littleEndian {
		msgId = binary.LittleEndian.Uint16(data[2:])
	} else {
		msgId = binary.BigEndian.Uint16(data[2:])
	}
	handler, ok := r.msgHandlers[msgId]
	if !ok {
		return msgId, errors.New("")
	}
	handler(msgId, data)
	return msgId, nil
}

func (r *Router) Register(msgId uint16, msgHandler MsgHandler) bool {
	if msgId > math.MaxUint16 {
		return false
	}
	if handler, ok := r.msgHandlers[msgId]; !ok {
		glog.Errorf("error", handler)
		return false
	}
	r.msgHandlers[msgId] = msgHandler
	return true
}

func (r *Router) Marshal(msgId uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, errors.New("")
	}
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		glog.Error(err)
		return data, err
	}
	buf := make([]byte, 4+len(data))
	if r.littleEndian {
		binary.LittleEndian.PutUint16(buf[0:2], 0)
		binary.LittleEndian.PutUint16(buf[2:], msgId)
	} else {

	}
	copy(buf[4:], data)
	return buf, nil
}

func (r *Router) UnMarshal(data []byte, msg interface{}) error {
	if len(data) < 4 {
		return errors.New("protobuf data too short")
	}
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("msg is not protobuf message")
	}
	return proto.Unmarshal(data[4:], pbMsg)
}
