package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableTeam(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_team",
		Description: "Teams are containers that contain scenarios and data accessible only by the members of the team.",
		List: &plugin.ListConfig{
			Hydrate: listTeams,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getTeam,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Team ID."},
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "The ID of the Organization."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Team."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")}},
	}
}

func getTeam(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getTeam", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var id = int(d.EqualsQuals["id"].GetInt64Value())
	var config = makesdk.NewRequestConfig(fmt.Sprintf(`teams/%d`, id))
	makesdk.ColumnsToParams(&config.Params, []string{"id", "name", "organizationId"})

	// fetch data
	var result = &makesdk.TeamResponse{}
	err = c.Get(config, &result)
	if err != nil {
		plugin.Logger(ctx).Error("make_team.getTeam", "request_error", err)
		return nil, c.HandleKnownErrors(err, "teams:read")
	}

	return result.Team, nil
}

func listTeams(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listTeams", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var op = makesdk.NewOrganizationListPaginator(c, -1)
	for op.HasMorePages() {
		organizations, err := op.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_team.listTeams", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			var up = makesdk.NewTeamListPaginator(c, int(d.RowsRemaining(ctx)), organization.Id)
			for up.HasMorePages() {
				teams, err := up.NextPage()
				if err != nil {
					plugin.Logger(ctx).Error("make_team.listTeams", "request_error", err)
					return nil, err
				}

				for _, i := range teams {
					d.StreamListItem(ctx, i)

					if d.RowsRemaining(ctx) <= 0 {
						return nil, nil
					}
				}
			}
		}
	}

	return nil, nil
}
