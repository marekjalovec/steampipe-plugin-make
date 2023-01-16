package make

import (
	"context"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Full name of the Role."},
			{Name: "subsidiary", Type: proto.ColumnType_BOOL, Description: "Is this Role defined in an Organization, or is it part of the account?"},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Can this role be used on the Organization, or Team level?"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Description: "Permissions of the users in this Role."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func listUserRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("listUserRoles", d, h)

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var config = client.NewRequestConfig("users/roles")
	utils.ColumnsToParams(&config.Params, []string{"id", "name", "subsidiary", "category", "permissions"})
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var pagesLeft = true
	for pagesLeft {
		var result = &client.UserRoleListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			logger.Error("make_user_role.listUserRoles", "request_error", err)
			return nil, c.HandleKnownErrors(err, "user:read")
		}

		// stream results
		for _, i := range result.UserRoles {
			d.StreamListItem(ctx, i)
		}

		// pagination
		var resultCount = len(result.UserRoles)
		if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
