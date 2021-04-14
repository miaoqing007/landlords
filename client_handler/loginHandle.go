package client_handler

import (
	"landlords/client_proto"
	"landlords/session"
)

func P_user_login_req(sess *session.Session, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_login_info(data)

	if tbl.F_account == "" || tbl.F_password == "" {
		return Code["error_ack"], client_proto.S_error_ack{"账号密码错误"}
	}

	if err := sess.InitPlayer(tbl.F_account, tbl.F_password); err != nil {
		return Code["error_ack"], client_proto.S_error_ack{"账号密码错误"}
	}

	return Code["user_login_req"], nil
}
