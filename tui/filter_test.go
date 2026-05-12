package tui_test

import (
	"testing"

	"github.com/anuress/miru/model"
	"github.com/anuress/miru/tui"
)

func req(method, url string, status int) model.Request {
	return model.Request{Method: method, URL: url, StatusCode: status}
}

func TestFilter_URLSubstring(t *testing.T) {
	f := tui.NewFilter("/api/auth")
	if !f.Match(req("GET", "https://example.com/api/auth/login", 200)) {
		t.Error("should match URL substring")
	}
	if f.Match(req("GET", "https://example.com/api/users", 200)) {
		t.Error("should not match")
	}
}

func TestFilter_Method(t *testing.T) {
	f := tui.NewFilter("m:POST")
	if !f.Match(req("POST", "https://x.com", 200)) {
		t.Error("should match POST")
	}
	if f.Match(req("GET", "https://x.com", 200)) {
		t.Error("should not match GET")
	}
}

func TestFilter_StatusPrefix(t *testing.T) {
	f := tui.NewFilter("s:4")
	if !f.Match(req("GET", "https://x.com", 401)) {
		t.Error("should match 401")
	}
	if !f.Match(req("GET", "https://x.com", 404)) {
		t.Error("should match 404")
	}
	if f.Match(req("GET", "https://x.com", 200)) {
		t.Error("should not match 200")
	}
}

func TestFilter_ExactStatus(t *testing.T) {
	f := tui.NewFilter("s:401")
	if !f.Match(req("GET", "https://x.com", 401)) {
		t.Error("should match 401")
	}
	if f.Match(req("GET", "https://x.com", 404)) {
		t.Error("should not match 404")
	}
}

func TestFilter_Empty(t *testing.T) {
	f := tui.NewFilter("")
	if !f.Match(req("GET", "https://x.com", 200)) {
		t.Error("empty filter should match all")
	}
}
