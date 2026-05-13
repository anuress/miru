package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anuress/miru/model"
)

func GenerateCurl(r model.Request) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "curl -X %s \\\n", r.Method)
	fmt.Fprintf(&sb, "  '%s'", r.URL)

	keys := make([]string, 0, len(r.ReqHeaders))
	for k := range r.ReqHeaders {
		keys = append(keys, k)
	}
	// Headers curl sets automatically — including them causes conflicts or unreadable output
	skipHeaders := map[string]bool{
		"content-length":   true,
		"host":             true,
		"connection":       true,
		"transfer-encoding": true,
		"accept-encoding":  true,
	}

	sort.Strings(keys)
	for _, k := range keys {
		if skipHeaders[strings.ToLower(k)] {
			continue
		}
		fmt.Fprintf(&sb, " \\\n  -H '%s: %s'", k, r.ReqHeaders[k])
	}

	if r.ReqBody != "" {
		fmt.Fprintf(&sb, " \\\n  -d '%s'", r.ReqBody)
	}
	return sb.String()
}
