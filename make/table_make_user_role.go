package make

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableUserRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_user_role",
		Description: "User roles usable for users within your Organizations and Teams. The column `category` signifies where the role can be used.",
		List: &plugin.ListConfig{
			Hydrate: listUserRoles,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The user ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Role."},
			{Name: "subsidiary", Type: proto.ColumnType_BOOL, Description: "Is the Role defined in an Organization, or is it part of the account?"},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Can the Role be used on the Organization, or Team level?"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Description: "Permissions of the users in the Role."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func listUserRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listUserRoles", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	urlp := c.NewUserRoleListPaginator(int(d.RowsRemaining(ctx)))
	for urlp.HasMorePages() {
		userRoles, err := urlp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_user_role.listUserRoles", "request_error", err)
			return nil, err
		}

		for _, i := range userRoles {
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
