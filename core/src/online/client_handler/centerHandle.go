package client_handler

import (
	"landlords/session"
)

//进入房间
func P_join_room_req(sess *session.Session, data []byte) (int16, interface{}) {
	//tbl, _ := client_proto.PKT_auto_id(data)
	//manager.AddPlayer2PvpPool(int(tbl.F_id), sess.User.Id, sess.User.Name)
	//return Code["join_room_req"], tbl
	return 0, nil
}
