package make

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableHook(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_hook",
		Description: "Webhooks allow you to send data to Make over HTTP. Webhooks create a URL that you can call from an external app or service, or from another Make scenario. Use webhooks to trigger the execution of scenarios.",
		List: &plugin.ListConfig{
			Hydrate: listHooks,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getHook,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Hook ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Hook."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "The Team ID."},
			{Name: "udid", Type: proto.ColumnType_STRING, Description: "The unique ID of the Hook."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the Hook (web/mail)."},
			{Name: "type_name", Type: proto.ColumnType_STRING, Description: "Full name of the Hook type."},
			{Name: "package_name", Type: proto.ColumnType_STRING, Description: "Package name."},
			{Name: "theme", Type: proto.ColumnType_STRING, Description: "The color theme of the Hook."},
			{Name: "flags", Type: proto.ColumnType_JSON, Description: "Additional flags of the Hook."},
			{Name: "is_editable", Type: proto.ColumnType_BOOL, Description: "Is the Hook editable?"},
			{Name: "is_enabled", Type: proto.ColumnType_BOOL, Description: "Is the Hook enabled?"},
			{Name: "is_gone", Type: proto.ColumnType_BOOL, Description: "Is the Hook gone?"},
			{Name: "queue_count", Type: proto.ColumnType_INT, Description: "Queue size for this Hook."},
			{Name: "queue_limit", Type: proto.ColumnType_INT, Description: "Queue limit for this Hook."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Additional metadata of the Hook."},
			{Name: "scenario_id", Type: proto.ColumnType_INT, Description: "The ID of the Scenario linked to this Hook."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "Full Hook address (URL, or email address)."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getHook(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getHook", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	hook, err := c.GetHook(int(d.EqualsQuals["id"].GetInt64Value()))
	if err != nil {
		plugin.Logger(ctx).Error("make_hook.getHook", "request_error", err)
		return nil, err
	}

	return hook, nil
}

func listHooks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listHooks", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_hook.listHooks", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				hlp := c.NewHookListPaginator(int(d.RowsRemaining(ctx)), team.Id)
				for hlp.HasMorePages() {
					hooks, err := hlp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_hook.listHooks", "request_error", err)
						return nil, err
					}

					for _, i := range hooks {
						i.TeamId = team.Id
						d.StreamListItem(ctx, i)

						if d.RowsRemaining(ctx) <= 0 {
							return nil, nil
						}
					}
				}
			}
		}
	}

	return nil, nil
}
