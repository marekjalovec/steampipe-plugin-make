package make

import (
	"context"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"net/url"
)

func ColumnsToParams(params *url.Values, columns []string) {
	for _, c := range columns {
		if c != "_ctx" {
			params.Add("cols", c)
		}
	}
}

func LogQueryContext(namespace string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) {
	plugin.Logger(ctx).Info(namespace, "Table", d.Table.Name)
	plugin.Logger(ctx).Info(namespace, "QueryContext", client.ToJSON(d.QueryContext))
	plugin.Logger(ctx).Info(namespace, "KeyColumnQuals", client.ToJSON(d.KeyColumnQuals))
	plugin.Logger(ctx).Info(namespace, "HydrateData", client.ToJSON(h))
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
