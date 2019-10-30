package client_handler

import (
	"app/misc/packet"
	"app/session"
)

func P_user_login_req(sess *session.Session, packet *packet.Packet) [][]byte {
	sess.InitUser()
	return nil
}
