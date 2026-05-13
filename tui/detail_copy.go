package tui

import (
	"regexp"
	"strings"
)

// jsonValueRe matches JSON key:value lines, capturing the value portion.
var jsonValueRe = regexp.MustCompile(`^\s*"[^"]*"\s*:\s*(.+?)[\s,]*$`)

// headerRe matches "Key: Value" header lines, capturing the value.
var headerRe = regexp.MustCompile(`^[^:]+:\s*(.+)$`)

// extractValue returns the value portion of a plain (ANSI-stripped) line.
// Returns "" if the line opens a JSON block — caller should use blockCopy instead.
func extractValue(line string) string {
	line = strings.TrimSpace(line)

	// JSON key:value pattern
	if m := jsonValueRe.FindStringSubmatch(line); len(m) == 2 {
		val := strings.TrimSpace(m[1])
		// String value — strip quotes
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			return val[1 : len(val)-1]
		}
		// Block opener — signal to caller
		if val == "{" || val == "[" || strings.HasSuffix(val, "{") || strings.HasSuffix(val, "[") {
			return ""
		}
		// Strip trailing comma
		val = strings.TrimRight(val, ",")
		return strings.TrimSpace(val)
	}

	// Header key: value pattern
	if m := headerRe.FindStringSubmatch(line); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}

	// Plain line — return trimmed
	return strings.TrimSpace(line)
}

// blockCopy collects lines from startIdx through the matching closing brace/bracket.
// Falls back to all lines from startIdx if no matching close is found.
func blockCopy(lines []string, startIdx int) string {
	if startIdx >= len(lines) {
		return ""
	}

	startLine := lines[startIdx]
	openIdx := strings.LastIndexAny(startLine, "{[")
	if openIdx == -1 {
		return strings.TrimSpace(startLine)
	}

	var open, close rune
	if startLine[openIdx] == '{' {
		open, close = '{', '}'
	} else {
		open, close = '[', ']'
	}

	depth := 0
	for i := startIdx; i < len(lines); i++ {
		for _, ch := range lines[i] {
			if ch == open {
				depth++
			} else if ch == close {
				depth--
				if depth == 0 {
					return strings.Join(lines[startIdx:i+1], "\n")
				}
			}
		}
	}
	return strings.Join(lines[startIdx:], "\n")
}

// isBlockOpener reports whether a plain line opens a JSON object or array.
func isBlockOpener(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasSuffix(trimmed, "{") || strings.HasSuffix(trimmed, "[")
}

// stripANSI removes ANSI escape sequences from s.
func stripANSI(s string) string {
	inEsc := false
	var b strings.Builder
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// stripANSILines strips ANSI codes from every line in the slice.
func stripANSILines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = stripANSI(l)
	}
	return out
}
