package make

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableScenarioLog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_scenario_log",
		Description: "Scenario Logs allow you to explore past runs of your scenarios.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("scenario_id"),
			Hydrate:    listScenarioLogs,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "scenario_id"}),
			Hydrate:    getScenarioLog,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The Scenario Log ID."},

			// Other Columns
			{Name: "imt_id", Type: proto.ColumnType_STRING, Description: "The Internal execution ID."},
			{Name: "duration", Type: proto.ColumnType_INT, Description: "Duration of the run in milliseconds."},
			{Name: "operations", Type: proto.ColumnType_INT, Description: "Number of operations consumed."},
			{Name: "transfer", Type: proto.ColumnType_INT, Description: "Transfer consumed in bytes."},
			{Name: "author_id", Type: proto.ColumnType_INT, Description: "Author User ID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the run."},
			{Name: "instant", Type: proto.ColumnType_BOOL, Description: "Instant run?"},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time of the run."},
			{Name: "status", Type: proto.ColumnType_INT, Description: "Status of the run."},
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "The Organization where the scenario belongs to."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "The Team where the scenario belongs to."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("ImtId")},

			// Virtual columns for the query
			{Name: "scenario_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func getScenarioLog(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getScenarioLog", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	scenarioId := int(d.EqualsQuals["scenario_id"].GetInt64Value())
	scenarioLogId := d.EqualsQuals["id"].GetStringValue()
	scenarioLog, err := c.GetScenarioLog(scenarioId, scenarioLogId)
	if err != nil {
		plugin.Logger(ctx).Error("make_scenario_log.getScenarioLog", "request_error", err)
		return nil, err
	}

	return scenarioLog, nil
}

func listScenarioLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listScenarioLogs", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	scenarioId := int(d.EqualsQuals["scenario_id"].GetInt64Value())
	sllp := c.NewScenarioLogListPaginator(int(d.RowsRemaining(ctx)), scenarioId)
	for sllp.HasMorePages() {
		scenarioLogs, err := sllp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_scenario_log.listScenarioLogs", "request_error", err)
			return nil, err
		}

		for _, i := range scenarioLogs {
			i.ScenarioId = scenarioId
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
