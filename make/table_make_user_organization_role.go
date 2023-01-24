package make

import (
	"context"
	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableUserOrganizationRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_user_organization_role",
		Description: "Users within your account.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("user_id"),
			Hydrate:    listUserOrganizationRoles,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The user ID."},

			// Other Columns
			{Name: "users_role_id", Type: proto.ColumnType_INT, Description: "The ID of the Role."},
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "The ID of the Organization."},
			{Name: "invitation", Type: proto.ColumnType_STRING, Description: "Is the invitation is still pending?"},
			{Name: "sso_pending", Type: proto.ColumnType_BOOL, Description: "Is SSO pending?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func listUserOrganizationRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listUserOrganizationRoles", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var userId = int(d.EqualsQuals["user_id"].GetInt64Value())
	var uorp = makesdk.NewUserOrganizationRoleListPaginator(c, int(d.RowsRemaining(ctx)), userId)
	for uorp.HasMorePages() {
		organizationRoles, err := uorp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_user_organization_role.listUserOrganizationRoles", "request_error", err)
			return nil, err
		}

		for _, i := range organizationRoles {
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
