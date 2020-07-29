package client_handler

import (
	"landlords/client_proto"
	"landlords/manager"
	"landlords/misc/packet"
	"landlords/wsconnection"
)

//进入房间
func P_join_room_req(ws *wsconnection.WsConnection, reader *packet.Packet) (int16, interface{}) {
	tbl, _ := client_proto.PKT_auto_id(reader)
	manager.AddPlayer2PvpPool(int(tbl.F_id), ws.User.Id)
	return nil
}
