package client_handler

import (
	"landlords/client_proto"
	"landlords/misc/packet"
	"landlords/model"
	"landlords/redis"
	"landlords/wsconnection"
)

func P_user_data_req(ws *wsconnection.WsConnection, reader *packet.Packet)(int16, interface{}) {
	if ws.User.Name == "" {
		return [][]byte{packet.Pack(Code["user_new_notify"], nil, nil)}
	}
	return nil
}

func P_user_reg_req(ws *wsconnection.WsConnection, reader *packet.Packet) (int16, interface{}) {
	tbl, _ := client_proto.PKT_entity_id(reader)
	if ws.User.Name != "" || tbl.F_id == "" {
		return nil
	}
	if redis.HExists(model.NAMEIDKEY, tbl.F_id) {
		return nil
	}
	ws.User.SetNameId(tbl.F_id)
	return nil
}
