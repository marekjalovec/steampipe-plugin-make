package make

import (
	"context"
	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOrganizationVariable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_organization_variable",
		Description: "Organization Variables are user-set variables you can use in your scenarios.",
		List: &plugin.ListConfig{
			Hydrate: listOrganizationVariables,
		},
		Columns: []*plugin.Column{
			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Organization Variable."},
			{Name: "type_id", Type: proto.ColumnType_INT, Description: "Original data type of the Organization Variable."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of the Organization Variable, converted to string."},
			{Name: "is_system", Type: proto.ColumnType_BOOL, Description: "Is the Organization Variable set by Make, or by users?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func listOrganizationVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listOrganizationVariables", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var op = makesdk.NewOrganizationListPaginator(c, -1)
	for op.HasMorePages() {
		organizations, err := op.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_organization_variable.listOrganizationVariables", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			var up = makesdk.NewOrganizationVariableListPaginator(c, int(d.RowsRemaining(ctx)), organization.Id)
			for up.HasMorePages() {
				organizationVariables, err := up.NextPage()
				if err != nil {
					plugin.Logger(ctx).Error("make_organization_variable.listOrganizationVariables", "request_error", err)
					return nil, err
				}

				for _, i := range organizationVariables {
					i.OrganizationId = organization.Id
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
