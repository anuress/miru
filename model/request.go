package model

import "time"

type Request struct {
	ID          string
	Method      string
	URL         string
	StatusCode  int
	Duration    time.Duration
	ReqHeaders  map[string]string
	ReqBody     string
	RespHeaders map[string]string
	RespBody    string
	RespBodyType string // populated from RespHeaders["Content-Type"]
	InFlight    bool
	Error       string
}
