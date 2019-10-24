package clienthandle

import (
	"app/misc/packet"
	"app/session"
	"fmt"
)

func P_heart_beat_req(session *session.Session, reader *packet.Packet) [][]byte {
	fmt.Println("heart")
	return [][]byte{
		packet.Pack(369, nil, nil),
	}
}
