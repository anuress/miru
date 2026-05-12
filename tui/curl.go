package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nurliman/miru/model"
)

func GenerateCurl(r model.Request) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "curl -X %s \\\n", r.Method)
	fmt.Fprintf(&sb, "  '%s'", r.URL)

	keys := make([]string, 0, len(r.ReqHeaders))
	for k := range r.ReqHeaders {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&sb, " \\\n  -H '%s: %s'", k, r.ReqHeaders[k])
	}

	if r.ReqBody != "" {
		fmt.Fprintf(&sb, " \\\n  -d '%s'", r.ReqBody)
	}
	return sb.String()
}
