package analytics

import "sort"

// TopK returns the top k words by count
func TopK[T comparable](dict map[T]int, top int) map[T]int {
	type kv struct {
		Key   T
		Value int
	}

	// clamp k to the number of unique words to avoid out of bounds
	top = min(top, len(dict))

	kvSlice := make([]kv, 0, len(dict))
	for k, v := range dict {
		kvSlice = append(kvSlice, kv{k, v})
	}

	sort.Slice(kvSlice, func(i, j int) bool {
		// descending order
		return kvSlice[i].Value > kvSlice[j].Value
	})

	result := make(map[T]int, top)
	for i := 0; i < top; i++ {
		result[kvSlice[i].Key] = kvSlice[i].Value
	}

	for i := top; i < len(kvSlice); i++ {
		// if there are more words with the same count as the kth word, include them
		if kvSlice[i].Value == kvSlice[top-1].Value {
			result[kvSlice[i].Key] = kvSlice[i].Value
		} else {
			break
		}
	}

	return result
}
