package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/client"
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
			{Name: "is_paused", Type: proto.ColumnType_BOOL, Description: "Is the organization paused?"},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "Make private instances use the externalId parameter for security reasons."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getOrganization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getOrganization", ctx, d, h)

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var id = int(d.EqualsQuals["id"].GetInt64Value())
	var config = client.NewRequestConfig(fmt.Sprintf(`organizations/%d`, id))
	ColumnsToParams(&config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})

	// fetch data
	var result = &client.OrganizationResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		plugin.Logger(ctx).Error("make_organization.getOrganization", "request_error", err)
		return nil, c.HandleKnownErrors(err, "organizations:read")
	}

	return result.Organization, nil
}

func listOrganizations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listOrganizations", ctx, d, h)

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var config = client.NewRequestConfig("organizations")
	ColumnsToParams(&config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var pagesLeft = true
	for pagesLeft {
		var result = &client.OrganizationListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			plugin.Logger(ctx).Error("make_organization.listOrganizations", "request_error", err)
			return nil, c.HandleKnownErrors(err, "organizations:read")
		}

		// stream results
		for _, i := range result.Organizations {
			d.StreamListItem(ctx, i)
		}

		// pagination
		var resultCount = len(result.Organizations)
		if d.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
