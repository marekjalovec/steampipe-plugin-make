package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableUserTeamRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_user_team_role",
		Description: "Users within your account.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("user_id"),
			Hydrate:    listUserTeamRoles,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The user ID."},

			// Other Columns
			{Name: "users_role_id", Type: proto.ColumnType_INT, Description: "ID of the Role."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "The ID of the Team."},
			{Name: "changeable", Type: proto.ColumnType_BOOL, Description: "Can this Role be changed?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func listUserTeamRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("listUserTeamRoles", d, h)

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var userId = int(d.KeyColumnQuals["user_id"].GetInt64Value())
	var config = client.NewRequestConfig(fmt.Sprintf(`users/%d/user-team-roles`, userId))
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var pagesLeft = true
	for pagesLeft {
		var result = &client.UserTeamRoleListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			logger.Error("make_user_team_role.listUserTeamRoles", "request_error", err)
			return nil, c.HandleKnownErrors(err, "user:read")
		}

		// stream results
		for _, i := range result.UserTeamRoles {
			d.StreamListItem(ctx, i)
		}

		// pagination
		var resultCount = len(result.UserTeamRoles)
		if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
