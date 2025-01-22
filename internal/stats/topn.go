package stats

import (
	"fmt"
	"maps"
	"slices"
)

func topNLanguages(n int, langBytes map[string]int, langColors map[string]string) []Language {
	totalBytes := 0
	for _, v := range langBytes {
		totalBytes += v
	}

	langs := slices.Collect(maps.Keys(langBytes))
	slices.SortFunc(langs, func(a string, b string) int {
		// desc. should return negative if b < a
		return langBytes[b] - langBytes[a]
	})

	// pick up to the top n languages, calculate their percentages, get their colors
	// and add to a slice
	out := make([]Language, 0, n)
	for i := 0; i < len(langs) && i < n; i++ {
		l := langs[i]
		lang := Language{
			Language: l,
			Percent:  fmt.Sprintf("%.1f", 100*float64(langBytes[l])/float64(totalBytes)),
			Color:    langColors[l],
		}
		out = append(out, lang)
	}

	return out
}
