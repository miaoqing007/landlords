package client_handle

import (
	"app/misc/packet"
	"app/session"
	"fmt"
)

func P_user_login_req(sess *session.Session, packet *packet.Packet) [][]byte {
	sess.InitUser()
	fmt.Println("sess", sess.UserName, sess.UserId)
	return nil
}
