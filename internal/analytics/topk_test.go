package analytics

import "testing"

func Test_TopK_Sanity(t *testing.T) {
	dict := map[string]int{
		"foo":  1,
		"bar":  2,
		"baz":  3,
		"quux": 4,
	}

	top := TopK(dict, 2)

	if len(top) != 2 {
		t.Errorf("expected 2 elements in top, got %d", len(top))
	}

	if top["quux"] != 4 {
		t.Errorf("expected quux to have count 4, got %d", top["quux"])
	}

	if top["baz"] != 3 {
		t.Errorf("expected baz to have count 3, got %d", top["baz"])
	}
}

func Test_TopK_Overflow(t *testing.T) {
	dict := map[string]int{
		"foo":  1,
		"bar":  2,
		"baz":  3,
		"quux": 4,
	}

	top := TopK(dict, 999)

	if len(top) != len(dict) {
		t.Errorf("expected same number of elements in dict and top, got %d", len(top))
	}

	for k, v := range top {
		if dict[k] != v {
			t.Errorf("expected %q to have count %d, got %d", k, dict[k], v)
		}
	}
}
