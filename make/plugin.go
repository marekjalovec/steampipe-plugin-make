/*
Package make implements a steampipe plugin for Make.com.

This plugin provides data that Steampipe uses to present foreign
tables that represent Make.com resources.
*/
package make

import (
	"context"
	"github.com/marekjalovec/steampipe-plugin-make/make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

const pluginName = "steampipe-plugin-make"

// Plugin creates this (make) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromJSONTag().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: client.ConfigInstance,
			Schema:      client.ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"make_organization": tableOrganization(ctx),
		},
	}

	utils.CreateLogger(ctx)

	return p
}
