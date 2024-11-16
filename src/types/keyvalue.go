package types

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// signature for map function
type MapFunc func(filename string, contents string) []KeyValue

// signature for reduce function
type ReduceFunc func(key string, values []string) string
