package mapfuncs

import (
	"mapreduce-go/src/types"
	"strconv"
	"strings"
)

func WordCountMapper(filename string, contents string) []types.KeyValue {
	// Split contents into words
	words := strings.Fields(contents)

	var kvs []types.KeyValue
	for _, word := range words {
		// remove punctuation, convert to lowercase
		cleanWord := strings.ToLower(strings.Trim(word, ".,!?;:\"'()[]{}"))
		if cleanWord != "" {
			kvs = append(kvs, types.KeyValue{Key: cleanWord, Value: "1"})
		}
	}

	return kvs
}

func WordCountReducer(key string, values []string) string {
	count := 0
	for _, value := range values {
		if num, err := strconv.Atoi(value); err == nil {
			count += num
		}
	}
	return strconv.Itoa(count)
}
