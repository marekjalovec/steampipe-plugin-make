package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strconv"
)

func tableTeam(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_team",
		Description: "Teams are containers that contain scenarios and data accessible only by the members of the team.",
		List: &plugin.ListConfig{
			Hydrate:       listTeams,
			ParentHydrate: listOrganizations,
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

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var id = int(d.KeyColumnQuals["id"].GetInt64Value())
	var config = client.NewRequestConfig(fmt.Sprintf(`teams/%d`, id))
	ColumnsToParams(&config.Params, []string{"id", "name", "organizationId"})

	// fetch data
	var result = &client.TeamResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		plugin.Logger(ctx).Error("make_team.getTeam", "request_error", err)
		return nil, c.HandleKnownErrors(err, "teams:read")
	}

	return result.Team, nil
}

func listTeams(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listTeams", ctx, d, h)

	if h.Item == nil {
		return nil, fmt.Errorf("no parent item found")
	}

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var config = client.NewRequestConfig("teams")
	ColumnsToParams(&config.Params, []string{"id", "name", "organizationId"})
	config.Params.Set("organizationId", strconv.Itoa(h.Item.(client.Organization).Id))
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var pagesLeft = true
	for pagesLeft {
		var result = &client.TeamListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			plugin.Logger(ctx).Error("make_team.listTeams", "request_error", err)
			return nil, c.HandleKnownErrors(err, "teams:read")
		}

		// stream results
		for _, i := range result.Teams {
			d.StreamListItem(ctx, i)
		}

		// pagination
		var resultCount = len(result.Teams)
		if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
