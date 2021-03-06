package gqldeduplicator

import (
	"encoding/json"
	"fmt"
)

// Deflate similar object in graphql response.
// Use deep first search (DFS) algorithm to walk over nodes and memoize object.
// If object appeared or memoized before, then it will deflated.
func Deflate(data []byte) ([]byte, error) {
	var node interface{}
	err := json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	memoize := make(map[string]bool)
	result := deflate(node, memoize, "root")

	return json.Marshal(result)
}

func deflate(node interface{}, memoize map[string]bool, path string) interface{} {
	switch value := node.(type) {
	case []interface{}:
		for i, v := range value {
			switch v.(type) {
			case []interface{}, map[string]interface{}:
				value[i] = deflate(v, memoize, path)
			default:
				value[i] = v
			}
		}
		return value
	case map[string]interface{}:
		if value != nil && value["id"] != nil && value["__typename"] != nil {
			key := fmt.Sprintf("%s,%v,%v", path, value["__typename"], value["id"])
			if memoize[key] {
				return map[string]interface{}{
					"id":         value["id"],
					"__typename": value["__typename"],
				}
			}

			memoize[key] = true
		}

		for k, v := range value {
			switch v.(type) {
			case []interface{}, map[string]interface{}:
				value[k] = deflate(v, memoize, path+","+k)
			default:
				value[k] = v
			}
		}

		return value
	}

	return node
}
