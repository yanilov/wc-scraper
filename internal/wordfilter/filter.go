package wordfilter

import (
	"unicode"

	"github.com/yanilov/wc-scraper/internal/bank"
)

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

func NewWordBankFilter(bank bank.WordBank) WordFilter {
	return func(word string) bool {
		return bank.Contains(word)
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
