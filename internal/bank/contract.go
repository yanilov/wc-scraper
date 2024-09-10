package bank

type WordBank interface {
	// Add adds a word to the bank
	Add(word string)

	// Contains checks if a word is in the bank
	Contains(word string) bool
}
