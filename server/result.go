package server

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"time"
)

type Result struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

func ResultByOk(data any) []byte {
	result := &Result{
		Code:      fasthttp.StatusOK,
		Message:   "ok",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
	bs, _ := json.Marshal(result)
	return bs
}

func NewResult(code int, message string, data any) []byte {
	result := &Result{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
	bs, _ := json.Marshal(result)
	return bs
}
