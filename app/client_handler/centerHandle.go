package client_handler

import (
	"landlords/app/client_proto"
	"landlords/app/manager"
	"landlords/app/misc/packet"
	"landlords/app/session"
)

//进入房间
func P_join_room_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_auto_id(reader)
	manager.AddPlayer2PvpPool(int(tbl.F_id), sess.User.Id)
	return nil
}
