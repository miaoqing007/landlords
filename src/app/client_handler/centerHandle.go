package client_handler

import (
	"app/client_proto"
	"app/manager"
	"app/misc/packet"
	"app/session"
)

//进入房间
func P_join_room_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_entity_id(reader)
	room := manager.GetRoomManager("1")
	if room == nil {
		return nil
	}
	if !room.AddPlayerToRoom(tbl.F_id) {
		return nil
	}
	return nil
}
