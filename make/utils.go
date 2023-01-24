package make

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func NewMakeClient(c *plugin.Connection) (*makesdk.Client, error) {
	var cfg, ok = c.Config.(Config)
	if !ok {
		return nil, fmt.Errorf("config object is not valid")
	}

	var config, err = makesdk.NewConfig(cfg.ApiToken, cfg.EnvironmentUrl, cfg.RateLimit)
	if err != nil {
		return nil, err
	}

	return makesdk.GetClient(config), nil
}

func ToJSON(value interface{}) string {
	j, _ := json.Marshal(value)
	return string(j)
}

func LogQueryContext(namespace string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) {
	plugin.Logger(ctx).Info(namespace, "Table", d.Table.Name)
	plugin.Logger(ctx).Info(namespace, "QueryContext", ToJSON(d.QueryContext))
	plugin.Logger(ctx).Info(namespace, "EqualsQuals", ToJSON(d.EqualsQuals))
	plugin.Logger(ctx).Info(namespace, "HydrateData", ToJSON(h))
}

func StandardColumnDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "The display name for the resource."
	case "virtual":
		return "Virtual column, used to map the entity to another object."
	}
	return ""
}
