package tui_test

import (
	"strings"
	"testing"

	"github.com/anuress/miru/model"
	"github.com/anuress/miru/tui"
)

func TestGenerateCurl_GET(t *testing.T) {
	r := model.Request{
		Method: "GET",
		URL:    "https://api.example.com/users",
		ReqHeaders: map[string]string{
			"Authorization": "Bearer tok",
		},
	}
	curl := tui.GenerateCurl(r)
	if !strings.Contains(curl, "curl -X GET") {
		t.Error("missing method")
	}
	if !strings.Contains(curl, "'https://api.example.com/users'") {
		t.Error("missing url")
	}
	if !strings.Contains(curl, "-H 'Authorization: Bearer tok'") {
		t.Error("missing header")
	}
}

func TestGenerateCurl_POST_WithBody(t *testing.T) {
	r := model.Request{
		Method:  "POST",
		URL:     "https://api.example.com/login",
		ReqBody: `{"email":"a@b.com"}`,
		ReqHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	curl := tui.GenerateCurl(r)
	if !strings.Contains(curl, "-d '{\"email\":\"a@b.com\"}'") {
		t.Error("missing body")
	}
}
