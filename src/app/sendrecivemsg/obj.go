package sendrecivemsg

var (
	SendMsgChan   = make(chan []byte, 1)
	ReciveMsgChan = make(chan []byte, 1)
)

type Resp struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
}

type Msg struct {
	Data string `json:"data"`
	Type int16  `json:"type"`
}
