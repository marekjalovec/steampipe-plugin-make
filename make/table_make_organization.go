package make

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_organization",
		Description: "Organizations are main containers that contain all teams, scenarios, and users.",
		List: &plugin.ListConfig{
			Hydrate: listOrganizations,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOrganization,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The organization ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Organization."},
			{Name: "country_id", Type: proto.ColumnType_INT, Description: "The ID of the country associated with the Organization."},
			{Name: "timezone_id", Type: proto.ColumnType_INT, Description: "The ID of the timezone associated with the Organization."},
			{Name: "license", Type: proto.ColumnType_JSON, Description: "Licence information and limits."},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "Zone where the Organization exists."},
			{Name: "service_name", Type: proto.ColumnType_STRING, Description: "Service name."},
			{Name: "is_paused", Type: proto.ColumnType_BOOL, Description: "Is the Organization paused?"},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "Make private instances use the externalId parameter for security reasons."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getOrganization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getOrganization", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	organization, err := c.GetOrganization(int(d.EqualsQuals["id"].GetInt64Value()))
	if err != nil {
		plugin.Logger(ctx).Error("make_organization.getOrganization", "request_error", err)
		return nil, err
	}

	return organization, nil
}

func listOrganizations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listOrganizations", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(int(d.RowsRemaining(ctx)))
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_organization.listOrganizations", "request_error", err)
			return nil, err
		}

		// stream results
		for _, i := range organizations {
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
