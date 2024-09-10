package wordfilter

import "unicode"

func NewMinLengthFilter(minLen int) WordFilter {
	return func(word string) bool {
		return len(word) >= minLen
	}
}

func NewAlphaOnlyFilter() WordFilter {
	return func(word string) bool {
		for _, char := range word {
			if !unicode.IsLetter(char) {
				return false
			}
		}
		return true
	}
}

func NewWordBankFilter(dict map[string]struct{}) WordFilter {
	return func(word string) bool {
		_, ok := dict[word]
		return ok
	}
}

func NewAggregateFilter(filters ...WordFilter) WordFilter {
	return func(word string) bool {
		for _, filter := range filters {
			if !filter(word) {
				return false
			}
		}
		return true
	}
}
