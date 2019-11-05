package client_handler

import (
	"app/client_proto"
	"app/misc/packet"
	"app/session"
)

func P_user_login_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_entity_id(reader)
	if err := sess.InitPlayer(tbl.F_id); err != nil {
		return [][]byte{
			packet.Pack(Code["error_ack"], nil, nil),
		}
	}
	return [][]byte{
		packet.Pack(Code["user_login_req"], tbl, nil),
	}
}
