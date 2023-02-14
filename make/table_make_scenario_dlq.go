package make

import (
	"context"
	makesdk "github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableScenarioDlq(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_scenario_dlq",
		Description: "If a scenario terminates unexpectedly because of an error, then the scenario run is discarded. You can set the scenario to store the failed scenario run as an incomplete execution. With that, if an error occurs in your scenario, you can resolve it manually and avoid losing data.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("scenario_id"),
			Hydrate:    listScenarioDlqs,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "scenario_id"}),
			Hydrate:    getScenarioDlq,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The Scenario Log ID."},

			// Other Columns
			{Name: "reason", Type: proto.ColumnType_STRING, Description: "The Internal execution ID."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "The Internal execution ID."},
			{Name: "index", Type: proto.ColumnType_INT, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
			{Name: "retry", Type: proto.ColumnType_BOOL, Description: "The Internal execution ID."},
			{Name: "attempts", Type: proto.ColumnType_INT, Description: "The Internal execution ID."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "The Internal execution ID."},
			{Name: "resolved", Type: proto.ColumnType_BOOL, Description: "The Internal execution ID."},
			{Name: "deleted", Type: proto.ColumnType_BOOL, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
			{Name: "execution_id", Type: proto.ColumnType_STRING, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
			{Name: "scenario_id", Type: proto.ColumnType_INT, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
			{Name: "scenario_name", Type: proto.ColumnType_STRING, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
			{Name: "team_name", Type: proto.ColumnType_STRING, Description: "The Internal execution ID.", Hydrate: getScenarioDlq},
		},
	}
}

func getScenarioDlq(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getScenarioDlq", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var id string
	if h.Item != nil {
		id = h.Item.(makesdk.ScenarioDlq).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}
	scenarioDlq, err := c.GetScenarioDlq(id)
	if err != nil {
		plugin.Logger(ctx).Error("make_scenario_log.getScenarioDlq", "request_error", err)
		return nil, err
	}

	return scenarioDlq, nil
}

func listScenarioDlqs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listScenarioDlqs", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	scenarioId := int(d.EqualsQuals["scenario_id"].GetInt64Value())
	up := c.NewScenarioDlqListPaginator(int(d.RowsRemaining(ctx)), scenarioId)
	for up.HasMorePages() {
		scenarioDlqs, err := up.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_scenario_log.listScenarioDlqs", "request_error", err)
			return nil, err
		}

		for _, i := range scenarioDlqs {
			i.ScenarioId = scenarioId
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
