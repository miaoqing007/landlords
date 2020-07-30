package client_handler

import (
	"landlords/client_proto"
	"landlords/initcards"
	"landlords/manager"
	"landlords/wsconnection"
)

//发牌
func P_licensing_card_req(ws *wsconnection.WsConnection, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_entity_id(data)
	info := client_proto.S_player_card{}
	cards := initcards.ShuffCards()
	room := manager.GetRoomManager(tbl.F_id)
	if room == nil {
		return Code["error_ack"], nil
	}
	room.CreatePlayerCards(cards[:17], cards[17:34], cards[34:51], cards[51:], &info)
	return Code["licensing_card_req"], info
}

//出牌
func P_out_of_the_card_req(ws *wsconnection.WsConnection, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_player_outof_card(data)
	if len(tbl.F_cards) == 0 {
		return Code["error_ack"], nil
	}
	room := manager.GetRoomManager(tbl.F_roomId)
	if room == nil {
		return Code["error_ack"], nil
	}
	if !room.CheckHandCards(ws.User.Id, tbl.F_cards) {
		return Code["error_ack"], nil
	}
	room.DeleteCards(ws.User.Id, tbl.F_cards)
	return Code["error_ack"], tbl
}
