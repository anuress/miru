package model

import "time"

type Request struct {
	ID           string
	Method       string
	URL          string
	StatusCode   int
	Duration     time.Duration
	ReqHeaders   map[string]string
	ReqBody      string
	ReqBodyType  string
	RespHeaders  map[string]string
	RespBody     string
	RespBodyType string
	StartedAt    time.Time
	InFlight     bool
}

func (r Request) StatusColour() string {
	switch {
	case r.InFlight:
		return "blue"
	case r.StatusCode >= 500:
		return "orange"
	case r.StatusCode >= 400:
		return "red"
	case r.StatusCode >= 300:
		return "orange"
	case r.StatusCode >= 200:
		return "green"
	default:
		return "white"
	}
}
