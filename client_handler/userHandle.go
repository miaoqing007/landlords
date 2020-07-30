package client_handler

import (
	"landlords/client_proto"
	"landlords/model"
	"landlords/redis"
	"landlords/wsconnection"
)

func P_user_data_req(ws *wsconnection.WsConnection, data []byte) (int16, interface{}) {
	if ws.User.Name == "" {
		return Code["error_ack"], nil
	}
	return Code["user_data_req"], nil
}

func P_user_reg_req(ws *wsconnection.WsConnection, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_entity_id(data)
	if ws.User.Name != "" || tbl.F_id == "" {
		return Code["error_ack"], nil
	}
	if redis.HExists(model.NAMEIDKEY, tbl.F_id) {
		return Code["error_ack"], nil
	}
	ws.User.SetNameId(tbl.F_id)
	return Code["user_reg_req"], nil
}
