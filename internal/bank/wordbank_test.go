package bank

import "testing"

func Test_Wordbank_New(t *testing.T) {
	words := []string{"hello", "world", "foo", "bar"}

	bank := New()
	for _, word := range words {
		bank.Add(word)
	}

	for _, word := range words {
		if !bank.Contains(word) {
			t.Errorf("expected %q to be in bank", word)
		}
	}
}

func Test_Wordbank_FromSlice(t *testing.T) {
	words := []string{"hello", "world", "foo", "bar"}

	bank := NewFromSlice(words)

	for _, word := range words {
		if !bank.Contains(word) {
			t.Errorf("expected %q to be in bank", word)
		}
	}
}

func Test_Wordbank_FromMap(t *testing.T) {
	words := []string{"hello", "world", "foo", "bar"}

	bank := NewFromSlice(words)

	for _, word := range words {
		if !bank.Contains(word) {
			t.Errorf("expected %q to be in bank", word)
		}
	}
}
