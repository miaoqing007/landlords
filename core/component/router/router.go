package router

import (
	"core/component/logger"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math"
)

type MsgHandler func(msgId uint16, data []byte)

type GatewayOnlineHandler func(msgId uint16, data []byte)

type OnlinePvpHandler func(msgId uint16, data []byte)

type Router struct {
	littleEndian         bool
	msgHandlers          map[uint16]MsgHandler           //玩家消息
	gatewayOnlineHandler map[uint16]GatewayOnlineHandler //gateway<-->online grpc消息
	onlinePvpHandler     map[uint16]OnlinePvpHandler     //online<-->pvp grpc消息
}

func NewRouter() *Router {
	r := &Router{
		littleEndian:         true,
		msgHandlers:          make(map[uint16]MsgHandler),
		gatewayOnlineHandler: make(map[uint16]GatewayOnlineHandler),
		onlinePvpHandler:     make(map[uint16]OnlinePvpHandler),
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

func (r *Router) RouterGatewayOnlineMsg(data []byte) (uint16, error) {
	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}
	var msgId uint16
	if r.littleEndian {
		msgId = binary.LittleEndian.Uint16(data[2:])
	} else {
		msgId = binary.BigEndian.Uint16(data[2:])
	}
	handler, ok := r.gatewayOnlineHandler[msgId]
	if !ok {
		return msgId, errors.New("")
	}
	handler(msgId, data)
	return msgId, nil
}

func (r *Router) RouterOnlinePvpMsg(data []byte) (uint16, error) {
	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}
	var msgId uint16
	if r.littleEndian {
		msgId = binary.LittleEndian.Uint16(data[2:])
	} else {
		msgId = binary.BigEndian.Uint16(data[2:])
	}
	handler, ok := r.onlinePvpHandler[msgId]
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
	if handler, ok := r.msgHandlers[msgId]; ok {
		logger.Errorf("error", handler)
		return false
	}
	r.msgHandlers[msgId] = msgHandler
	return true
}

func (r *Router) RegisterGatewayOnline(msgId uint16, msgHandler GatewayOnlineHandler) bool {
	if handler, ok := r.gatewayOnlineHandler[msgId]; ok {
		logger.Errorf("error", handler)
		return false
	}
	r.gatewayOnlineHandler[msgId] = msgHandler
	return true
}

func (r *Router) RegisterOnlinePvp(msgId uint16, msgHandler OnlinePvpHandler) bool {
	if handler, ok := r.onlinePvpHandler[msgId]; ok {
		logger.Errorf("error", handler)
		return false
	}
	r.onlinePvpHandler[msgId] = msgHandler
	return true
}

func (r *Router) Marshal(msgId uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, errors.New("")
	}
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		logger.Error(err)
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