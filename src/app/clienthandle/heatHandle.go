package clienthandle

import (
	"app/misc/packet"
	"app/session"
	"fmt"
)

func P_heart_beat_req(session *session.Session, packet *packet.Packet) [][]byte {
	fmt.Println("heart")

	return nil
}
