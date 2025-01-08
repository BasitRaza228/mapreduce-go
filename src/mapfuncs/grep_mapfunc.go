package mapfuncs

import (
	"fmt"
	"mapreduce-go/src/types"
	"strings"
)

func GrepMapper(pattern string) types.MapFunc {
	return func(filename string, contents string) []types.KeyValue {
		var kvs []types.KeyValue
		lines := strings.Split(contents, "\n")

		for lineNum, line := range lines {
			if strings.Contains(line, pattern) {
				key := fmt.Sprintf("%s:%d", filename, lineNum+1)
				kvs = append(kvs, types.KeyValue{Key: key, Value: line})
			}
		}

		return kvs
	}
}

// reduce function for grep-like functionality
func GrepReducer(key string, values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}
