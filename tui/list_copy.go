package tui

import (
	"sort"
	"strings"

	"github.com/anuress/miru/model"
)

// RespBodyCopy returns the pretty-printed response body for clipboard copy.
// Returns "" if the body is empty.
func RespBodyCopy(r model.Request) string {
	return prettyBody(r.RespBody)
}

// RawRequestCopy returns a plain-text representation of request headers and
// body for clipboard copy. Headers are sorted alphabetically. Returns "" if
// both headers and body are empty.
func RawRequestCopy(r model.Request) string {
	var parts []string

	if len(r.ReqHeaders) > 0 {
		keys := make([]string, 0, len(r.ReqHeaders))
		for k := range r.ReqHeaders {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var hdr strings.Builder
		for _, k := range keys {
			hdr.WriteString(k + ": " + r.ReqHeaders[k] + "\n")
		}
		parts = append(parts, strings.TrimRight(hdr.String(), "\n"))
	}

	if body := prettyBody(r.ReqBody); body != "" {
		parts = append(parts, body)
	}

	return strings.Join(parts, "\n")
}
