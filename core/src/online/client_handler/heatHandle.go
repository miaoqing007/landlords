package client_handler

import (
	"landlords/client_proto"
	"landlords/session"
)

//心跳检测
func P_heart_beat_req(sess *session.Session, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_auto_id(data)
	return Code["heart_beat_req"], tbl
}
