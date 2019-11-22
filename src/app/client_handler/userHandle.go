package client_handler

import (
	"app/client_proto"
	"app/misc/packet"
	"app/model"
	"app/redis"
	"app/session"
)

func P_user_data_req(sess *session.Session, reader *packet.Packet) [][]byte {
	if sess.User.Name == "" {
		return [][]byte{packet.Pack(Code["user_new_notify"], nil, nil)}
	}
	return nil
}

func P_user_reg_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_entity_id(reader)
	if sess.User.Name != "" || tbl.F_id == "" {
		return nil
	}
	if redis.HExists(model.NAMEIDKEY, tbl.F_id) {
		return nil
	}
	sess.User.SetNameId(tbl.F_id)
	return nil
}
