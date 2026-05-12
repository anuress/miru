package tui_test

import (
	"testing"

	"github.com/anuress/miru/tui"
)

func TestSearch_FindMatches(t *testing.T) {
	s := tui.NewSearch("title")
	body := `{"title":"Hello","id":1,"title2":"World"}`
	matches := s.FindMatches(body)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].Start != 2 {
		t.Errorf("wrong start: %d", matches[0].Start)
	}
}

func TestSearch_Empty(t *testing.T) {
	s := tui.NewSearch("")
	matches := s.FindMatches("anything")
	if len(matches) != 0 {
		t.Error("empty query should return no matches")
	}
}
