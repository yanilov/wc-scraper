package wordfilter

import (
	"testing"

	"github.com/yanilov/wc-scraper/internal/bank"
)

func Test_BankFilter(t *testing.T) {
	dict := map[string]struct{}{
		"hello": {},
		"world": {},
	}
	bank := bank.NewFromKeys(dict)
	valid := NewWordBankFilter(bank)
	if !valid("hello") {
		t.Errorf("expected 'hello' to pass filter")
	}
	if !valid("world") {
		t.Errorf("expected 'world' to pass filter")
	}
	if valid("foo") {
		t.Errorf("expected 'foo' to fail filter")
	}
}

func Test_MinLengthFilter(t *testing.T) {
	valid := NewMinLengthFilter(3)
	if !valid("foo") {
		t.Errorf("expected 'foo' to pass filter")
	}
	if !valid("foobar") {
		t.Errorf("expected 'foobar' to pass filter")
	}
	if valid("fo") {
		t.Errorf("expected 'fo' to fail filter")
	}
}

func Test_AlphaOnlyFilter(t *testing.T) {
	valid := NewAlphaOnlyFilter()
	if !valid("foo") {
		t.Errorf("expected 'foo' to pass filter")
	}
	if !valid("foobar") {
		t.Errorf("expected 'foobar' to pass filter")
	}
	if valid("foo1") {
		t.Errorf("expected 'foo1' to fail filter")
	}
}

func Test_AggregateFilter(t *testing.T) {

	dict := map[string]struct{}{
		"foo": {},
		"fo":  {},
		"bar": {},
	}
	bank := bank.NewFromKeys(dict)

	valid := NewAggregateFilter(
		NewMinLengthFilter(3),
		NewAlphaOnlyFilter(),
		NewWordBankFilter(bank),
	)
	if !valid("foo") {
		t.Errorf("expected 'foo' to pass filter")
	}
	if valid("foobar") {
		t.Errorf("expected 'foobar' to fail filter")
	}
	if valid("foo1") {
		t.Errorf("expected 'foo1' to fail filter")
	}
	if valid("fo") {
		t.Errorf("expected 'fo' to fail filter")
	}
}
