package client_handler

import (
	"app/client_proto"
	"app/helper/conv"
	"app/manager"
	"app/misc/packet"
	"app/session"
)

var id int

func P_join_room_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_entity_id(reader)
	room := manager.GetRoomManager(tbl.F_id)
	if !room.AddPlayerToRoom(conv.FormatInt(id)) {
		return nil
	}
	id++
	return nil
}
