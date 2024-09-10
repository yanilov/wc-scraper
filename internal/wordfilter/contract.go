package wordfilter

// WordFilter is a function that filters words. It returns true if the word should be included.
type WordFilter func(string) bool
