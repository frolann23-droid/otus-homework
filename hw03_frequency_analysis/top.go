package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type frequency struct {
	word  string
	count int
}

var withoutPunctRegexp = regexp.MustCompile(`^[[:punct:]]+|[[:punct:]]+$`)

func Top10(s string) []string {
	words := strings.Fields(s)

	wordsFrequency := make(map[string]*frequency)
	var frequencies []*frequency

	result := make([]string, 0, 10)

	for _, word := range words {
		cleanWord := strings.ToLower(withoutPunctRegexp.ReplaceAllString(word, ""))

		if cleanWord == "" {
			continue
		}

		freq, ok := wordsFrequency[cleanWord]

		if !ok {
			freq = &frequency{word: cleanWord, count: 0}
			wordsFrequency[cleanWord] = freq
			frequencies = append(frequencies, freq)
		}

		freq.count++
	}

	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}

		return frequencies[i].count > frequencies[j].count
	})

	for i, freq := range frequencies {
		if i > 9 {
			break
		}

		result = append(result, freq.word)
	}

	return result
}
