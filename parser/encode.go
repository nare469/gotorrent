package parser

import (
	"bytes"
	"sort"
	"strconv"
)

func encode(it Item) string {
	var buffer bytes.Buffer

	switch it.(type) {
	case map[string]Item:
		buffer.WriteString("d")

		var keys []string
		for k := range it.(map[string]Item) {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			buffer.WriteString(encode(k))
			buffer.WriteString(encode(it.(map[string]Item)[k]))
		}

		buffer.WriteString("e")
	case []Item:
		buffer.WriteString("l")

		for _, s := range it.([]Item) {
			buffer.WriteString(encode(s))
		}
	case int64:
		buffer.WriteString("i")
		buffer.WriteString(strconv.FormatInt(it.(int64), 10))
		buffer.WriteString("e")
	case string:
		buffer.WriteString(strconv.Itoa(len(it.(string))))
		buffer.WriteString(":")
		buffer.WriteString(it.(string))
	}

	return buffer.String()
}
