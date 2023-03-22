package make

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableTeamVariable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_team_variable",
		Description: "Team Variables are user-set variables you can use in your scenarios.",
		List: &plugin.ListConfig{
			Hydrate: listTeamVariables,
		},
		Columns: []*plugin.Column{
			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Team Variable."},
			{Name: "type_id", Type: proto.ColumnType_INT, Description: "Original data type of the Team Variable."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of the Team Variable, converted to string."},
			{Name: "is_system", Type: proto.ColumnType_BOOL, Description: "Is the Team Variable set by Make, or by users?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "team_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func listTeamVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listTeamVariables", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_team_variable.listTeamVariables", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				tvlp := c.NewTeamVariableListPaginator(int(d.RowsRemaining(ctx)), team.Id)
				for tvlp.HasMorePages() {
					teamVariables, err := tvlp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_team_variable.listTeamVariables", "request_error", err)
						return nil, err
					}

					for _, i := range teamVariables {
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
