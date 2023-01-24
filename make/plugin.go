/*
Package make implements a steampipe plugin for Make.com.

This plugin provides data that Steampipe uses to present foreign
tables that represent Make.com resources.
*/
package make

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-make"

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"make_api_token":              tableApiToken(ctx),
			"make_connection":             tableConnection(ctx),
			"make_data_store":             tableDataStore(ctx),
			"make_organization":           tableOrganization(ctx),
			"make_organization_variable":  tableOrganizationVariable(ctx),
			"make_team":                   tableTeam(ctx),
			"make_team_variable":          tableTeamVariable(ctx),
			"make_user":                   tableUser(ctx),
			"make_user_organization_role": tableUserOrganizationRole(ctx),
			"make_user_role":              tableUserRole(ctx),
			"make_user_team_role":         tableUserTeamRole(ctx),
		},
	}

	return p
}
