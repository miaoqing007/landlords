package client_handler

import (
	"github.com/golang/glog"
	"landlords/client_proto"
	"landlords/wsconnection"
)

func P_user_login_req(ws *wsconnection.WsConnection, data []byte) (int16, interface{}) {
	tbl, _ := client_proto.PKT_entity_id(data)
	if err := ws.InitPlayer(tbl.F_id); err != nil {
		return Code["error_ack"], nil
	}
	tbl.F_id = ws.User.Id
	glog.Info(tbl.F_id)
	return Code["user_login_req"], tbl
}
