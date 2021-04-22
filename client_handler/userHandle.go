package client_handler

import (
	"landlords/client_proto"
	"landlords/model"
	"landlords/redis"
	"landlords/session"
)

func P_user_data_req(sess *session.Session, data []byte) (int16, interface{}) {
	return Code["user_data_req"], client_proto.S_user_info{sess.User.Name, sess.User.Id}
}

func P_user_reg_req(sess *session.Session, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_entity_id(data)
	if sess.User.Name != "" || tbl.F_id == "" {
		return Code["error_ack"], nil
	}
	if redis.HExists(model.NAMEIDKEY, tbl.F_id) {
		return Code["error_ack"], nil
	}
	sess.User.SetNameId(tbl.F_id)
	return Code["user_reg_req"], nil
}
