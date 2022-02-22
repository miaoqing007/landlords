package network

import (
	"errors"
	"fmt"
	"io"
	"math"
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	lenMsgLen int32  // 需要读取的包长度
	minMsgLen uint32 // 最小长度 包体包含了cmdID 所以最小长度至少为2
	maxMsgLen uint32 // 允许的最大的包体长度
	recvBuff  *ByteBuffer
	sendBuff  *ByteBuffer
}

func newMsgParser() *MsgParser {
	msgParser := &MsgParser{
		lenMsgLen: 4,
		minMsgLen: 2,
		maxMsgLen: 2 * 1024 * 1024,
		recvBuff:  NewByteBuffer(),
		sendBuff:  NewByteBuffer(),
	}
	return msgParser
}

// SetMsgLen It's dangerous to call the method on reading or writing
func (p *MsgParser) SetMsgLen(lenMsgLen int32, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}

	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// Read goroutine safe
func (p *MsgParser) Read(conn *TcpConn) (uint16, []byte, error) {

	p.recvBuff.EnsureWritableBytes(p.lenMsgLen)

	readLen, err := io.ReadFull(conn, p.recvBuff.WriteBuff()[:p.lenMsgLen])
	// read len
	if err != nil {
		return 0, nil, fmt.Errorf("%v readLen:%v", err, readLen)
	}
	p.recvBuff.WriteBytes(int32(readLen))

	// parse len
	var msgLen uint32
	switch p.lenMsgLen {
	case 2:
		msgLen = uint32(p.recvBuff.ReadInt16())
	case 4:
		msgLen = uint32(p.recvBuff.ReadInt32())
	}

	// check len
	if msgLen > p.maxMsgLen {
		//logger.Errorf("message too long msgLen %d maxMsgLen %d", msgLen, p.maxMsgLen)
		return 0, nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		//logger.Errorf("message too short msgLen %d minMsgLen %d", msgLen, p.minMsgLen)
		return 0, nil, errors.New("message too short")
	}

	p.recvBuff.EnsureWritableBytes(int32(msgLen))

	rLen, err := io.ReadFull(conn, p.recvBuff.WriteBuff()[:msgLen])
	if err != nil {
		return 0, nil, fmt.Errorf("%v msgLen:%v readLen:%v", err, msgLen, rLen)
	}
	p.recvBuff.WriteBytes(int32(rLen))

	// 保留了2字节flag 暂时未处理
	var flag uint16

	flag = uint16(p.recvBuff.ReadInt16())

	// 减去2字节的保留字段长度
	return flag, p.recvBuff.NextBytes(int32(msgLen - 2)), nil

}

// goroutine safe
func (p *MsgParser) Write(conn *TcpConn, buff ...byte) error {
	// get len
	msgLen := uint32(len(buff))

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	// write len
	switch p.lenMsgLen {
	case 2:
		p.sendBuff.AppendInt16(int16(msgLen))
	case 4:
		p.sendBuff.AppendInt32(int32(msgLen))
	}

	p.sendBuff.Append(buff)
	// write data
	writeBuff := p.sendBuff.ReadBuff()[:p.sendBuff.Length()]

	_, err := conn.Write(writeBuff)

	p.sendBuff.Reset()

	return err
}

func (p *MsgParser) reset() {
	p.recvBuff = NewByteBuffer()
	p.sendBuff = NewByteBuffer()
}
