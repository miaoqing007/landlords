package client_handler

import (
	"landlords/client_proto"
	"landlords/misc/packet"
	"landlords/session"
)

//心跳检测
func P_heart_beat_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_auto_id(reader)
	return [][]byte{
		packet.Pack(Code["heart_beat_req"], tbl, nil),
	}
}
