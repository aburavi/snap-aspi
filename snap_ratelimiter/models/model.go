package models

type Trx struct {
	TrxType          string `json:"trx_type"`
	SystemTraceAudit string `json:"system_trace_audit"`
	TrxDateTime      string `json:"trx_date_time"`
	ClientID         string `json:"client_id"`
	PosTerminalType  string `json:"pos_terminal_type"`
	Sign             string `json:"sign"`
	AccountNo        string `json:"account_no"`
}

type StreamData struct {
	TrxType          string `json:"trx_type"`
	SystemTraceAudit string `json:"system_trace_audit"`
	Data             string `json:"data"`
	UUID             string `json:"uuid"`
	IData            interface{}
}

type Worker struct {
	Object chan interface{}
	Quit   chan struct{}
}

type UUIDObject struct {
	UUID       string
	StreamData chan StreamData
	Command    string
}
