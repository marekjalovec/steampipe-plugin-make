package utils

import (
	"encoding/json"
	"net/url"
)

func ToJSON(value interface{}) string {
	j, _ := json.Marshal(value)
	return string(j)
}

func ColumnsToParams(params *url.Values, columns []string) {
	for _, c := range columns {
		if c != "_ctx" {
			params.Add("cols", c)
		}
	}
}
