package protocol_test

import (
	"strings"
	"testing"

	"github.com/nurliman/miru/protocol"
)

func TestReadMessage_Request(t *testing.T) {
	line := `{"id":"abc","type":"request","method":"GET","url":"https://api.example.com/users","headers":{"Authorization":"Bearer tok"},"body":""}`
	msg, err := protocol.ParseMessage([]byte(line))
	if err != nil {
		t.Fatal(err)
	}
	if msg.ID != "abc" {
		t.Errorf("bad id: %s", msg.ID)
	}
	if msg.Method != "GET" {
		t.Errorf("bad method: %s", msg.Method)
	}
	if msg.URL != "https://api.example.com/users" {
		t.Errorf("bad url: %s", msg.URL)
	}
}

func TestReadMessage_Response(t *testing.T) {
	line := `{"id":"abc","type":"response","status":200,"duration":142,"headers":{"Content-Type":"application/json"},"body":"{\"name\":\"Nur\"}"}`
	msg, err := protocol.ParseMessage([]byte(line))
	if err != nil {
		t.Fatal(err)
	}
	if msg.Status != 200 {
		t.Errorf("bad status: %d", msg.Status)
	}
	if msg.Duration != 142 {
		t.Errorf("bad duration: %d", msg.Duration)
	}
}

func TestReader_MultipleMessages(t *testing.T) {
	data := strings.NewReader(`{"id":"1","type":"request","method":"POST","url":"https://x.com/login","headers":{},"body":"{}"}
{"id":"1","type":"response","status":401,"duration":89,"headers":{},"body":"{}"}
`)
	msgs, err := protocol.ReadAll(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
}
