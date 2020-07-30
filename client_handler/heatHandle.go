package client_handler

import (
	"landlords/client_proto"
	"landlords/wsconnection"
)

//心跳检测
func P_heart_beat_req(ws *wsconnection.WsConnection, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_auto_id(data)
	return Code["heart_beat_req"], tbl
}
