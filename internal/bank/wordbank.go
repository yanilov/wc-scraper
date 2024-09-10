package bank

type wordBankImpl struct {
	// holds a set of words which are part of the bank
	data map[string]struct{}
}

// assert that wordBankImpl implements WordBank
var _ WordBank = (*wordBankImpl)(nil)

// Add implements WordBank.
func (w *wordBankImpl) Add(word string) {
	panic("unimplemented")
}

// Contains implements WordBank.
func (w *wordBankImpl) Contains(word string) bool {
	panic("unimplemented")
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
	return NewFromMap(data)
}

func NewFromMap(words map[string]struct{}) WordBank {
	return &wordBankImpl{
		data: words,
	}
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
