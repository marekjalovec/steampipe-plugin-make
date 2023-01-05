package utils

import (
	"encoding/json"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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

func LogQueryContext(namespace string, d *plugin.QueryData, h *plugin.HydrateData) {
	logger.Info(namespace, "Table", d.Table.Name)
	logger.Info(namespace, "QueryContext", ToJSON(d.QueryContext))
	logger.Info(namespace, "KeyColumnQuals", ToJSON(d.KeyColumnQuals))
	logger.Info(namespace, "HydrateData", ToJSON(h))
}

func StandardColumnDescription(key string) string {
	switch key {
	case "title":
		return "The display name for this resource."
	case "virtual":
		return "Virtual column, used to map this entity to another object."
	}
	return ""
}
