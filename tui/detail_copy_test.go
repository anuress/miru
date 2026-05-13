package tui

import "testing"

func TestExtractValue_HeaderLine(t *testing.T) {
	got := extractValue("Authorization: Bearer eyJhbGci...")
	if got != "Bearer eyJhbGci..." {
		t.Errorf("got %q", got)
	}
}

func TestExtractValue_JSONString(t *testing.T) {
	got := extractValue(`  "name": "Nur",`)
	if got != "Nur" {
		t.Errorf("got %q", got)
	}
}

func TestExtractValue_JSONNumber(t *testing.T) {
	got := extractValue(`  "id": 42,`)
	if got != "42" {
		t.Errorf("got %q", got)
	}
}

func TestExtractValue_JSONBool(t *testing.T) {
	got := extractValue(`  "active": true,`)
	if got != "true" {
		t.Errorf("got %q", got)
	}
}

func TestExtractValue_JSONNull(t *testing.T) {
	got := extractValue(`  "token": null,`)
	if got != "null" {
		t.Errorf("got %q", got)
	}
}

func TestExtractValue_JSONObjectOpener(t *testing.T) {
	got := extractValue(`  "user": {`)
	if got != "" {
		t.Errorf("expected empty for object opener, got %q", got)
	}
}

func TestExtractValue_PlainLine(t *testing.T) {
	got := extractValue("  hello world  ")
	if got != "hello world" {
		t.Errorf("got %q", got)
	}
}

func TestBlockCopy_Object(t *testing.T) {
	lines := []string{
		`  "user": {`,
		`    "name": "Nur",`,
		`    "id": 42`,
		`  },`,
		`  "other": 1`,
	}
	got := blockCopy(lines, 0)
	want := "  \"user\": {\n    \"name\": \"Nur\",\n    \"id\": 42\n  },"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestBlockCopy_Array(t *testing.T) {
	lines := []string{
		`  "items": [`,
		`    1,`,
		`    2`,
		`  ],`,
	}
	got := blockCopy(lines, 0)
	want := "  \"items\": [\n    1,\n    2\n  ],"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestBlockCopy_UnmatchedBrace(t *testing.T) {
	lines := []string{`  "x": {`, `    "y": 1`}
	got := blockCopy(lines, 0)
	if got == "" {
		t.Error("expected non-empty fallback")
	}
}
