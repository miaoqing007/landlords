package client_handler

import (
	"landlords/client_proto"
	"landlords/model"
	"landlords/redis"
	"landlords/session"
)

func P_user_data_req(sess *session.Session, data []byte) (int16, interface{}) {
	info := client_proto.S_user_info{}
	if sess.User.Name == "" {
		return Code["register_name_ack"], nil
	}
	info.F_name = sess.User.Name
	info.F_uid = sess.User.Id
	return Code["user_data_req"], info
}

func P_register_name_req(sess *session.Session, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_msg_string(data)
	info := client_proto.S_user_info{}
	sess.User.SetName(tbl.F_msg)
	info.F_name = sess.User.Name
	info.F_uid = sess.User.Id
	return Code["user_data_req"], info
}

func P_user_reg_req(sess *session.Session, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_entity_id(data)
	if sess.User.Name != "" || tbl.F_id == "" {
		return Code["error_ack"], nil
	}
	if redis.HExists(model.NAMEIDKEY, tbl.F_id) {
		return Code["error_ack"], nil
	}
	return Code["user_reg_req"], nil
}
