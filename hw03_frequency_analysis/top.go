package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const topSize = 10

type kv struct {
	key   string
	value uint
}

var reg = regexp.MustCompile(`[^a-zA-Z0-9_а-яА-Я\-]+`)

func split(s string) []string {
	// return strings.Fields(input) without asterisk
	return reg.Split(strings.ReplaceAll(s, " - ", ""), -1)
}

func convert(s string) string {
	// return s without asterisk
	return strings.ToLower(s)
}

func Top10(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	words := split(input)

	counters := make(map[string]uint)
	for _, word := range words {
		counters[convert(word)]++
	}

	size := len(counters)
	pairs := make([]kv, size)
	i := 0
	for word, count := range counters {
		pairs[i] = kv{word, count}
		i++
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].value != pairs[j].value {
			return pairs[i].value > pairs[j].value
		}
		// lexicographically
		return pairs[i].key < pairs[j].key
	})

	// size = min(topSize, size) in go 1.21 or later
	if size > topSize {
		size = topSize
	}

	result := make([]string, size)
	for i := 0; i < size; i++ {
		result[i] = pairs[i].key
	}

	return result
}
