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
			{Name: "service_name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "is_paused", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "Make private instances use the externalId parameter for security reasons."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getOrganization(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("getOrganization", d, h)

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var id = int(d.KeyColumnQuals["id"].GetInt64Value())
	var config = client.NewRequestConfig(fmt.Sprintf(`organizations/%d`, id))
	utils.ColumnsToParams(&config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})

	// fetch data
	var result = &client.OrganizationResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		logger.Error("make_organization.getOrganization", "request_error", err)
		return nil, c.HandleKnownErrors(err, "organizations:read")
	}

	return result.Organization, nil
}

func listOrganizations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("listOrganizations", d, h)

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var config = client.NewRequestConfig("organizations")
	utils.ColumnsToParams(&config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var pagesLeft = true
	for pagesLeft {
		var result = &client.OrganizationListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			logger.Error("make_organization.listOrganizations", "request_error", err)
			return nil, c.HandleKnownErrors(err, "organizations:read")
		}

		// stream results
		for _, i := range result.Organizations {
			d.StreamListItem(ctx, i)
		}

		// pagination
		var resultCount = len(result.Organizations)
		if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
