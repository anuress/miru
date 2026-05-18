package tui_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anuress/miru/model"
	"github.com/anuress/miru/tui"
)

func addReqs(m *tui.ListModel, urls ...string) {
	for _, u := range urls {
		m.AddRequest(model.Request{URL: u})
	}
}

func key(s string) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func arrowKey(s string) tea.Msg {
	switch s {
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

// TestListNav_FilterThenNavigate verifies cursor moves through filtered results.
func TestListNav_FilterThenNavigate(t *testing.T) {
	m := tui.NewListModel()
	// Add 5 requests; 3 match "/api"
	addReqs(&m,
		"https://example.com/api/c",
		"https://example.com/other",
		"https://example.com/api/b",
		"https://example.com/health",
		"https://example.com/api/a",
	)

	m.SetFilter(tui.NewFilter("/api"))

	first := m.Selected()
	if first == nil {
		t.Fatal("expected selection after filter, got nil")
	}

	// Navigate down
	m, _ = m.Update(arrowKey("down"))
	second := m.Selected()
	if second == nil {
		t.Fatal("expected selection after down, got nil")
	}
	if second.URL == first.URL {
		t.Errorf("cursor did not advance: still at %s", first.URL)
	}

	// Navigate down again
	m, _ = m.Update(arrowKey("down"))
	third := m.Selected()
	if third == nil {
		t.Fatal("expected selection after second down, got nil")
	}
	if third.URL == second.URL {
		t.Errorf("cursor did not advance on second down: still at %s", second.URL)
	}

	// Should be at bottom now — one more down should stay
	m, _ = m.Update(arrowKey("down"))
	bottom := m.Selected()
	if bottom == nil || bottom.URL != third.URL {
		t.Errorf("expected cursor to stay at bottom, got %v", bottom)
	}
}

// TestListNav_NonMatchingRequestDoesNotShiftCursor verifies that an incoming request
// which doesn't match the active filter does not push the cursor out of visible range.
func TestListNav_NonMatchingRequestDoesNotShiftCursor(t *testing.T) {
	m := tui.NewListModel()
	addReqs(&m,
		"https://example.com/api/b",
		"https://example.com/api/a",
	)
	m.SetFilter(tui.NewFilter("/api"))

	// Navigate down to second item
	m, _ = m.Update(arrowKey("down"))
	second := m.Selected()
	if second == nil {
		t.Fatal("expected selection after down")
	}

	// Non-matching request arrives — cursor must not move
	m.AddRequest(model.Request{URL: "https://example.com/health"})
	still := m.Selected()
	if still == nil {
		t.Fatal("cursor became nil after non-matching AddRequest")
	}
	if still.URL != second.URL {
		t.Errorf("cursor shifted: expected %s, got %s", second.URL, still.URL)
	}

	// Navigation must still work after the non-matching request
	m, _ = m.Update(arrowKey("up"))
	if m.Selected() == nil {
		t.Fatal("navigation broken after non-matching AddRequest")
	}
}

// TestListNav_CursorResetOnFilter verifies that changing the filter resets cursor to 0.
func TestListNav_CursorResetOnFilter(t *testing.T) {
	m := tui.NewListModel()
	addReqs(&m,
		"https://example.com/api/a",
		"https://example.com/api/b",
		"https://example.com/api/c",
	)

	// Move cursor to position 2
	m, _ = m.Update(arrowKey("down"))
	m, _ = m.Update(arrowKey("down"))
	if m.Selected() == nil {
		t.Fatal("expected selection at position 2")
	}

	// Apply filter — cursor should reset
	m.SetFilter(tui.NewFilter("/api"))
	sel := m.Selected()
	if sel == nil {
		t.Fatal("expected selection after filter reset, got nil")
	}
}
