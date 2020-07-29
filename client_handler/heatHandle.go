package client_handler

import (
	"landlords/client_proto"
	"landlords/misc/packet"
	"landlords/wsconnection"
)

//心跳检测
func P_heart_beat_req(ws *wsconnection.WsConnection, reader *packet.Packet) (int16, interface{}) {
	tbl, _ := client_proto.PKT_auto_id(reader)
	return Code["heart_beat_req"], tbl
}
