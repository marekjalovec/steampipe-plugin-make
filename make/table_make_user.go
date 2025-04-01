package make

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_user",
		Description: "Existing Users within your account, and their attributes.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The user ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the User."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "User's email."},
			{Name: "language", Type: proto.ColumnType_STRING, Description: "Environment language."},
			{Name: "locale_id", Type: proto.ColumnType_INT, Description: "Environment locale - date and hour formats, decimal separators etc."},
			{Name: "country_id", Type: proto.ColumnType_INT, Description: "The ID of the country associated with the User."},
			{Name: "timezone_id", Type: proto.ColumnType_INT, Description: "The ID of the timezone associated with the User."},
			{Name: "features", Type: proto.ColumnType_JSON, Description: "Features enabled for the User."},
			{Name: "avatar", Type: proto.ColumnType_STRING, Description: "Gravatar URL for the User."},
			{Name: "last_login", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time of the last login."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listUsers", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_user.listUsers", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				ulp := c.NewUserListPaginator(-1, team.Id)
				for ulp.HasMorePages() {
					users, err := ulp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_user.listUsers", "request_error", err)
						return nil, err
					}

					for _, i := range users {
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
