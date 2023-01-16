package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strconv"
)

func tableUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_user",
		Description: "Existing Users within your account, and their attributes.",
		List: &plugin.ListConfig{
			Hydrate:       listUsers,
			ParentHydrate: listOrganizations,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The user ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Full name of the User."},
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

	if h.Item == nil {
		return nil, fmt.Errorf("parent organization not defined")
	}

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// iterate over organization teams
	var organization = h.Item.(client.Organization)
	var teams = organization.Teams
	for _, team := range teams {
		// prepare params
		var config = client.NewRequestConfig("users")
		config.Params.Set("teamId", strconv.Itoa(team.Id))
		if d.QueryContext.Limit != nil {
			config.Pagination.Limit = int(*d.QueryContext.Limit)
		}

		// fetch data
		var pagesLeft = true
		for pagesLeft {
			var result = &client.UserListResponse{}
			err = c.Get(&config, result)
			if err != nil {
				plugin.Logger(ctx).Error("make_user.listUsers", "request_error", err)
				return nil, c.HandleKnownErrors(err, "user:read")
			}

			// stream results
			for _, i := range result.Users {
				i.OrganizationId = organization.Id
				i.TeamId = team.Id
				d.StreamListItem(ctx, i)
			}

			// pagination
			var resultCount = len(result.Users)
			if d.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
				pagesLeft = false
			} else {
				config.Pagination.Offset += config.Pagination.Limit
			}
		}
	}

	return nil, nil
}
