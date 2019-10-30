package client_handler

import (
	"app/client_proto"
	"app/misc/packet"
	"app/session"
)

func P_heart_beat_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_auto_id(reader)
	return [][]byte{
		packet.Pack(Code["heart_beat_req"], tbl, nil),
	}
}
