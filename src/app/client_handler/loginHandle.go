package client_handler

import (
	"app/client_proto"
	"app/misc/packet"
	"app/session"
)

func P_user_login_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_entity_id(reader)
	sess.InitPlayer(tbl.F_id)
	return [][]byte{
		packet.Pack(Code["user_login_req"], tbl, nil),
	}
}
