package make

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			{Name: "changeable", Type: proto.ColumnType_BOOL, Description: "Can the Role be changed?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func listUserTeamRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listUserTeamRoles", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	userId := int(d.EqualsQuals["user_id"].GetInt64Value())
	utrlp := c.NewUserTeamRoleListPaginator(int(d.RowsRemaining(ctx)), userId)
	for utrlp.HasMorePages() {
		teamRoles, err := utrlp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_user_team_role.listUserTeamRoles", "request_error", err)
			return nil, err
		}

		for _, i := range teamRoles {
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
