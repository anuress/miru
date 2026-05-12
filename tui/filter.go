package tui

import (
	"fmt"
	"strings"

	"github.com/nurliman/miru/model"
)

type Filter struct {
	raw string
}

func NewFilter(input string) Filter {
	return Filter{raw: strings.TrimSpace(input)}
}

func (f Filter) Match(r model.Request) bool {
	if f.raw == "" {
		return true
	}
	if strings.HasPrefix(f.raw, "m:") {
		method := strings.ToUpper(strings.TrimPrefix(f.raw, "m:"))
		return strings.ToUpper(r.Method) == method
	}
	if strings.HasPrefix(f.raw, "s:") {
		prefix := strings.TrimPrefix(f.raw, "s:")
		return strings.HasPrefix(fmt.Sprintf("%d", r.StatusCode), prefix)
	}
	return strings.Contains(r.URL, f.raw)
}
