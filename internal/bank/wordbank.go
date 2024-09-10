package bank

type wordBankImpl struct {
	// holds a set of words which are part of the bank
	data map[string]struct{}
}

// assert that wordBankImpl implements WordBank
var _ WordBank = (*wordBankImpl)(nil)

// Add implements WordBank.
func (w *wordBankImpl) Add(word string) {
	w.data[word] = struct{}{}
}

// Contains implements WordBank.
func (w *wordBankImpl) Contains(word string) bool {
	_, ok := w.data[word]
	return ok
}

func New() WordBank {
	return &wordBankImpl{
		data: make(map[string]struct{}),
	}
}

func NewFromSlice(words []string) WordBank {
	data := make(map[string]struct{}, len(words))
	for _, word := range words {
		data[word] = struct{}{}
	}
	return NewFromKeys(data)
}

func NewFromKeys[V any](words map[string]V) WordBank {
	data := make(map[string]struct{}, len(words))
	for key := range words {
		data[key] = struct{}{}
	}
	return &wordBankImpl{
		data,
	}
}

func NewFromValues[K comparable](words map[K]string) WordBank {
	data := make(map[string]struct{}, len(words))
	for _, value := range words {
		data[value] = struct{}{}
	}
	return &wordBankImpl{
		data,
	}
}
