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
		// SchemaMode: plugin.SchemaModeDynamic,
		// DefaultTransform: transform.FromGo(),
		// DefaultGetConfig: &plugin.GetConfig{
		//	IgnoreConfig: &plugin.IgnoreConfig{
		//		ShouldIgnoreErrorFunc: isNotFoundError([]string{
		//			"NoSuchEntity",
		//			"NotFoundException",
		//			"ResourceNotFoundException",
		//			"InvalidParameter",
		//			"InvalidParameterValue",
		//			"InvalidParameterValueException",
		//			"ValidationError",
		//			"ValidationException",
		//		}),
		//	},
		// },
		// Default ignore config for the plugin
		// DefaultIgnoreConfig: &plugin.IgnoreConfig{
		//	ShouldIgnoreErrorFunc: shouldIgnoreErrorPluginDefault(),
		// },
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"make_api_token":              tableApiToken(ctx),
			"make_connection":             tableConnection(ctx),
			"make_data_store":             tableDataStore(ctx),
			"make_function":               tableFunctions(ctx),
			"make_organization":           tableOrganization(ctx),
			"make_organization_variable":  tableOrganizationVariable(ctx),
			"make_scenario":               tableScenario(ctx),
			"make_scenario_dlq":           tableScenarioDlq(ctx),
			"make_scenario_log":           tableScenarioLog(ctx),
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
