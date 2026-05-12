package protocol

import (
	"bufio"
	"encoding/json"
	"io"
)

type Message struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Status   int               `json:"status"`
	Duration int               `json:"duration"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

func ParseMessage(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return m, err
}

func ReadAll(r io.Reader) ([]Message, error) {
	var msgs []Message
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		m, err := ParseMessage(line)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, scanner.Err()
}

// NewStreamReader emits Messages as they arrive on r, closing the channel when r closes.
func NewStreamReader(r io.Reader) <-chan Message {
	ch := make(chan Message)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			m, err := ParseMessage(line)
			if err == nil {
				ch <- m
			}
		}
	}()
	return ch
}
