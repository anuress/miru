package tui

import (
	"strings"
	"testing"

	"github.com/anuress/miru/model"
)

func TestRespBodyCopy_JSON(t *testing.T) {
	r := model.Request{RespBody: `{"id":1,"name":"Alice"}`}
	got := RespBodyCopy(r)
	want := "{\n  \"id\": 1,\n  \"name\": \"Alice\"\n}"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRespBodyCopy_Empty(t *testing.T) {
	r := model.Request{}
	got := RespBodyCopy(r)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestRespBodyCopy_NonJSON(t *testing.T) {
	r := model.Request{RespBody: "OK"}
	got := RespBodyCopy(r)
	if got != "OK" {
		t.Errorf("got %q", got)
	}
}

func TestRawRequestCopy_HeadersAndBody(t *testing.T) {
	r := model.Request{
		ReqHeaders: map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		},
		ReqBody: `{"key":"value"}`,
	}
	got := RawRequestCopy(r)
	if !strings.Contains(got, "Authorization: Bearer token123") {
		t.Errorf("missing Authorization header:\n%s", got)
	}
	if !strings.Contains(got, "Content-Type: application/json") {
		t.Errorf("missing Content-Type header:\n%s", got)
	}
	if !strings.Contains(got, `"key": "value"`) {
		t.Errorf("missing pretty-printed body:\n%s", got)
	}
}

func TestRawRequestCopy_NoHeaders(t *testing.T) {
	r := model.Request{ReqBody: `{"x":1}`}
	got := RawRequestCopy(r)
	if got == "" {
		t.Error("expected non-empty result when body is present")
	}
	if strings.Contains(got, "no headers") {
		t.Errorf("must not contain placeholder text, got:\n%s", got)
	}
}

func TestRawRequestCopy_HeadersOnly(t *testing.T) {
	r := model.Request{
		ReqHeaders: map[string]string{"Accept": "application/json"},
	}
	got := RawRequestCopy(r)
	if !strings.Contains(got, "Accept: application/json") {
		t.Errorf("expected header line, got:\n%s", got)
	}
}

func TestRawRequestCopy_Empty(t *testing.T) {
	r := model.Request{}
	got := RawRequestCopy(r)
	if got != "" {
		t.Errorf("expected empty for empty request, got %q", got)
	}
}
