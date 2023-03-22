package make

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableScenario(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_scenario",
		Description: "Scenarios allow you to create and run automation tasks. A scenario consists of a series of modules that indicate how data should be transferred and transformed between apps or services.",
		List: &plugin.ListConfig{
			Hydrate: listScenarios,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getScenario,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Scenario ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Scenario."},
			{Name: "hook_id", Type: proto.ColumnType_INT, Description: "The Webhook ID, if used to trigger the Scenario."},
			{Name: "device_id", Type: proto.ColumnType_INT, Description: "The device ID."},
			{Name: "device_scope", Type: proto.ColumnType_STRING, Description: "The device scope."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The Scenario description."},
			{Name: "folder_id", Type: proto.ColumnType_INT, Description: "The folder id, if Scenario is in a folder."},
			{Name: "is_invalid", Type: proto.ColumnType_BOOL, Description: "Is the Scenario invalid?"},
			{Name: "is_linked", Type: proto.ColumnType_BOOL, Description: "Is the Scenario linked?"},
			{Name: "is_locked", Type: proto.ColumnType_BOOL, Description: "Is the Scenario locked?"},
			{Name: "is_paused", Type: proto.ColumnType_BOOL, Description: "Is the Scenario paused?"},
			{Name: "concept", Type: proto.ColumnType_BOOL, Description: "Is the Scenario a concept?"},
			{Name: "used_packages", Type: proto.ColumnType_JSON, Description: "The list of used Apps in this Scenario."},
			{Name: "last_edit", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time of when was this Scenario edited the last time."},
			{Name: "scheduling", Type: proto.ColumnType_JSON, Description: "Scheduling details"},
			{Name: "is_waiting", Type: proto.ColumnType_BOOL, Description: "Is the Scenario paused due to an error?"},
			{Name: "dlq_count", Type: proto.ColumnType_INT, Description: "Number of incomplete executions."},
			{Name: "created_by_user", Type: proto.ColumnType_JSON, Description: "Author of this Scenario."},
			{Name: "updated_by_user", Type: proto.ColumnType_JSON, Description: "Last editor of this Scenario."},
			{Name: "next_exec", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time of the next scheduled execution."},
			{Name: "scenario_version", Type: proto.ColumnType_INT, Description: "The version of the Scenario."},
			{Name: "module_sequence_id", Type: proto.ColumnType_INT, Description: "The module sequence ID."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func getScenario(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getScenario", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	scenario, err := c.GetScenario(int(d.EqualsQuals["id"].GetInt64Value()))
	if err != nil {
		plugin.Logger(ctx).Error("make_scenario.getScenario", "request_error", err)
		return nil, err
	}

	return scenario, nil
}

func listScenarios(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listScenarios", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_scenario.listScenarios", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				slp := c.NewScenarioListPaginator(int(d.RowsRemaining(ctx)), team.Id, 0)
				for slp.HasMorePages() {
					scenarios, err := slp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_scenario.listScenarios", "request_error", err)
						return nil, err
					}

					for _, i := range scenarios {
						i.OrganizationId = organization.Id
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
