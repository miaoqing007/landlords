package client_handler

import (
	"app/client_proto"
	"app/initcards"
	"app/manager"
	"app/misc/packet"
	"app/session"
)

//发牌
func P_licensing_card_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_entity_id(reader)
	info := client_proto.S_player_card{}
	cards := initcards.ShuffCards()
	room := manager.GetRoomManager(tbl.F_id)
	if room == nil {
		return nil
	}
	room.CreatePlayerCards(cards[:17], cards[17:34], cards[34:51], cards[51:], &info)
	return nil
}

func P_out_of_the_card_req(sess *session.Session, reader *packet.Packet) [][]byte {
	tbl, _ := client_proto.PKT_player_outof_card(reader)
	if len(tbl.F_cards) == 0 {
		return nil
	}
	room := manager.GetRoomManager(tbl.F_roomId)
	if room == nil {
		return nil
	}
	if !room.CheckHandCards(sess.User.Id, tbl.F_cards) {
		return nil
	}
	room.DeleteCards(sess.User.Id, tbl.F_cards)
	return nil
}
