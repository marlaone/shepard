package hashmap

type Pair[K comparable, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}
