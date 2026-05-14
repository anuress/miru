package protocol

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// Message is a fully assembled request+response reconstructed from OKPRFL logcat lines.
type Message struct {
	ID          string
	Method      string
	URL         string
	Duration    int // milliseconds, from RST
	Status      int // HTTP status code, from RSS
	ReqHeaders  map[string]string
	ReqBody     string
	RespHeaders map[string]string
	RespBody    string
	Error       string // set on REE (network error)
}

// NewStreamReader reads OKPRFL logcat lines from r, assembles them into complete
// Messages (one per request), and emits each on the returned channel.
// The channel is closed when r is closed or returns an error.
func NewStreamReader(r io.Reader) <-chan Message {
	ch := make(chan Message)
	go func() {
		defer close(ch)
		partial := make(map[string]*Message)
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 4*1024*1024), 4*1024*1024) // 4MB — handles large response bodies
		for scanner.Scan() {
			tag, value := ExtractOKPRFL(scanner.Text())
			if tag == "" {
				continue
			}
			id, typeCode := ParseTag(tag)
			if id == "" {
				continue
			}

			msg, exists := partial[id]
			if !exists {
				msg = &Message{
					ID:          id,
					ReqHeaders:  make(map[string]string),
					RespHeaders: make(map[string]string),
				}
				partial[id] = msg
			}

			switch typeCode {
			case "RQM":
				msg.Method = value
			case "RQU":
				msg.URL = value
			case "RQH":
				k, v := splitHeader(value)
				msg.ReqHeaders[k] = v
			case "RQB":
				msg.ReqBody += value
			case "RSS":
				msg.Status, _ = strconv.Atoi(value)
			case "RST":
				msg.Duration, _ = strconv.Atoi(value)
			case "RSH":
				k, v := splitHeader(value)
				msg.RespHeaders[k] = v
			case "RSB":
				msg.RespBody += value
			case "RSD":
				ch <- *msg
				delete(partial, id)
			case "REE":
				msg.Error = value
				ch <- *msg
				delete(partial, id)
			}
		}
	}()
	return ch
}

// ExtractOKPRFL finds the OKPRFL tag and value within a logcat line.
// Logcat line format: "MM-DD HH:MM:SS.mmm  PID  TID  LEVEL OKPRFL_<id>_<type>: value"
func ExtractOKPRFL(line string) (tag, value string) {
	idx := strings.Index(line, "OKPRFL_")
	if idx == -1 {
		return "", ""
	}
	rest := line[idx:]
	colon := strings.Index(rest, ": ")
	if colon == -1 {
		return "", ""
	}
	return rest[:colon], rest[colon+2:]
}

// ParseTag splits "OKPRFL_<id>_<typeCode>" into id and typeCode.
// typeCode is always the 3-char suffix after the last underscore.
func ParseTag(tag string) (id, typeCode string) {
	inner := strings.TrimPrefix(tag, "OKPRFL_")
	last := strings.LastIndex(inner, "_")
	if last == -1 {
		return "", ""
	}
	return inner[:last], inner[last+1:]
}

// splitHeader splits "Name: value" into key and value on first ": ".
func splitHeader(s string) (key, value string) {
	idx := strings.Index(s, ": ")
	if idx == -1 {
		return strings.TrimSpace(s), ""
	}
	return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+2:])
}
