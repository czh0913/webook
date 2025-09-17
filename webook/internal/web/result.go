package web

type Result struct {
	Code int    `json:"Code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
