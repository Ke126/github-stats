package stats

import (
	"fmt"
	"sort"
)

func topNLanguages(n int, langBytes map[string]int, langColors map[string]string) []Language {
	totalBytes := 0
	i := 0

	type MapElem struct {
		k string
		v int
	}
	s := make([]MapElem, len(langBytes))
	for k, v := range langBytes {
		totalBytes += v
		s[i] = MapElem{
			k: k,
			v: v,
		}
		i++
	}

	// sort in descending order
	sort.Slice(s, func(i, j int) bool {
		return s[i].v > s[j].v
	})

	// pick up to the top n languages, calculate their percentages, get their colors
	// and add to a slice
	out := make([]Language, 0, n)
	for i := 0; i < len(s) && i < n; i++ {
		lang := Language{
			Language: s[i].k,
			Percent:  fmt.Sprintf("%.1f", 100*float64(s[i].v)/float64(totalBytes)),
			Color:    langColors[s[i].k],
		}
		out = append(out, lang)
	}

	return out
}
