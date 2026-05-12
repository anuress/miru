package tui

import "strings"

type Match struct {
	Start int
	End   int
}

type Search struct {
	Query string
}

func NewSearch(query string) Search {
	return Search{Query: query}
}

func (s Search) FindMatches(text string) []Match {
	if s.Query == "" {
		return nil
	}
	var matches []Match
	lower := strings.ToLower(text)
	q := strings.ToLower(s.Query)
	offset := 0
	for {
		idx := strings.Index(lower[offset:], q)
		if idx == -1 {
			break
		}
		start := offset + idx
		end := start + len(q)
		matches = append(matches, Match{Start: start, End: end})
		offset = end
	}
	return matches
}
