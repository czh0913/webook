package web

type Result struct {
	code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
