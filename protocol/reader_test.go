package protocol_test

import (
	"strings"
	"testing"

	"github.com/anuress/miru/protocol"
)

func logLine(tag, value string) string {
	return "05-12 10:00:00.000  1234  5678 V " + tag + ": " + value
}

func TestExtractOKPRFL_ValidLine(t *testing.T) {
	line := logLine("OKPRFL_abc123_RQM", "GET")
	tag, value := protocol.ExtractOKPRFL(line)
	if tag != "OKPRFL_abc123_RQM" {
		t.Errorf("wrong tag: %q", tag)
	}
	if value != "GET" {
		t.Errorf("wrong value: %q", value)
	}
}

func TestExtractOKPRFL_NonOKPRFLLine(t *testing.T) {
	line := logLine("okhttp.OkHttpClient", "--> GET https://example.com")
	tag, value := protocol.ExtractOKPRFL(line)
	if tag != "" || value != "" {
		t.Errorf("expected empty, got tag=%q value=%q", tag, value)
	}
}

func TestParseTag(t *testing.T) {
	id, typeCode := protocol.ParseTag("OKPRFL_req-uuid-123_RQM")
	if id != "req-uuid-123" {
		t.Errorf("wrong id: %q", id)
	}
	if typeCode != "RQM" {
		t.Errorf("wrong typeCode: %q", typeCode)
	}
}

func TestNewStreamReader_CompleteRequest(t *testing.T) {
	lines := strings.Join([]string{
		logLine("OKPRFL_r1_RQM", "GET"),
		logLine("OKPRFL_r1_RQU", "https://api.example.com/users"),
		logLine("OKPRFL_r1_RQH", "Authorization: Bearer tok"),
		logLine("OKPRFL_r1_RSS", "200"),
		logLine("OKPRFL_r1_RST", "142"),
		logLine("OKPRFL_r1_RSH", "Content-Type: application/json"),
		logLine("OKPRFL_r1_RSB", `{"id":1}`),
		logLine("OKPRFL_r1_RSD", "-->"),
	}, "\n")

	ch := protocol.NewStreamReader(strings.NewReader(lines))
	msg := <-ch
	if msg.Method != "GET" {
		t.Errorf("wrong method: %q", msg.Method)
	}
	if msg.URL != "https://api.example.com/users" {
		t.Errorf("wrong url: %q", msg.URL)
	}
	if msg.Status != 200 {
		t.Errorf("wrong status: %d", msg.Status)
	}
	if msg.Duration != 142 {
		t.Errorf("wrong duration: %d", msg.Duration)
	}
	if msg.ReqHeaders["Authorization"] != "Bearer tok" {
		t.Errorf("wrong req header: %v", msg.ReqHeaders)
	}
	if msg.RespHeaders["Content-Type"] != "application/json" {
		t.Errorf("wrong resp header: %v", msg.RespHeaders)
	}
	if msg.RespBody != `{"id":1}` {
		t.Errorf("wrong body: %q", msg.RespBody)
	}
}

func TestNewStreamReader_ErrorRequest(t *testing.T) {
	lines := strings.Join([]string{
		logLine("OKPRFL_r2_RQM", "GET"),
		logLine("OKPRFL_r2_RQU", "https://api.example.com/fail"),
		logLine("OKPRFL_r2_REE", "java.net.UnknownHostException: Unable to resolve host"),
	}, "\n")

	ch := protocol.NewStreamReader(strings.NewReader(lines))
	msg := <-ch
	if msg.Error == "" {
		t.Error("expected error to be set")
	}
	if msg.Method != "GET" {
		t.Errorf("wrong method: %q", msg.Method)
	}
}

func TestNewStreamReader_InterleavedRequests(t *testing.T) {
	// Two requests interleaved by thread
	lines := strings.Join([]string{
		logLine("OKPRFL_r1_RQM", "GET"),
		logLine("OKPRFL_r2_RQM", "POST"),
		logLine("OKPRFL_r1_RQU", "https://api.example.com/a"),
		logLine("OKPRFL_r2_RQU", "https://api.example.com/b"),
		logLine("OKPRFL_r1_RSS", "200"),
		logLine("OKPRFL_r2_RSS", "201"),
		logLine("OKPRFL_r1_RST", "100"),
		logLine("OKPRFL_r2_RST", "200"),
		logLine("OKPRFL_r1_RSD", "-->"),
		logLine("OKPRFL_r2_RSD", "-->"),
	}, "\n")

	ch := protocol.NewStreamReader(strings.NewReader(lines))
	msgs := map[string]protocol.Message{}
	for m := range ch {
		msgs[m.Method] = m
	}
	if msgs["GET"].URL != "https://api.example.com/a" {
		t.Errorf("wrong GET url: %q", msgs["GET"].URL)
	}
	if msgs["POST"].URL != "https://api.example.com/b" {
		t.Errorf("wrong POST url: %q", msgs["POST"].URL)
	}
	if msgs["GET"].Status != 200 {
		t.Errorf("wrong GET status: %d", msgs["GET"].Status)
	}
	if msgs["POST"].Status != 201 {
		t.Errorf("wrong POST status: %d", msgs["POST"].Status)
	}
}

func TestNewStreamReader_ChunkedBody(t *testing.T) {
	lines := strings.Join([]string{
		logLine("OKPRFL_r1_RQM", "POST"),
		logLine("OKPRFL_r1_RQU", "https://api.example.com/upload"),
		logLine("OKPRFL_r1_RSB", `{"part":1}`),
		logLine("OKPRFL_r1_RSB", `{"part":2}`),
		logLine("OKPRFL_r1_RSS", "200"),
		logLine("OKPRFL_r1_RST", "300"),
		logLine("OKPRFL_r1_RSD", "-->"),
	}, "\n")

	ch := protocol.NewStreamReader(strings.NewReader(lines))
	msg := <-ch
	if msg.RespBody != `{"part":1}{"part":2}` {
		t.Errorf("chunked body not concatenated: %q", msg.RespBody)
	}
}
