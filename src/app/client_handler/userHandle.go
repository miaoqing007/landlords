package client_handler

import (
	"app/misc/packet"
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
	if redis.Exists("") {
		return nil
	}
	return nil
}
