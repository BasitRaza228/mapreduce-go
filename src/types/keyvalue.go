package types

// represents Map of MapReduce Framework
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// signature for map functions
type MapFunc func(filename string, contents string) []KeyValue

// signature for reduce functions
type ReduceFunc func(key string, values []string) string
