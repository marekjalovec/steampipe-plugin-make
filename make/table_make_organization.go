package make

import (
	"context"
	"github.com/marekjalovec/steampipe-plugin-make/make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type Organization struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	CountryId   int                 `json:"countryId"`
	TimezoneId  int                 `json:"timezoneId"`
	License     OrganizationLicence `json:"license"`
	Zone        string              `json:"zone"`
	ServiceName string              `json:"serviceName"`
	IsPaused    bool                `json:"isPaused"`
	ExternalId  string              `json:"externalId"`
	Teams       []OrganizationTeam  `json:"teams"` // used to load make_connection
}

type OrganizationTeam struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type OrganizationLicence struct {
	Apps       []string `json:"apps"`
	Users      int      `json:"users"`
	Dslimit    int64    `json:"dslimit"`
	Fslimit    int64    `json:"fslimit"`
	Iolimit    int64    `json:"iolimit"`
	Dsslimit   int64    `json:"dsslimit"`
	Fulltext   bool     `json:"fulltext"`
	Interval   int      `json:"interval"`
	Transfer   int64    `json:"transfer"`
	Operations int64    `json:"operations"`
}

type OrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type OrganizationListResponse struct {
	Organizations []Organization `json:"organizations"`
	Pagination    Pagination     `json:"pg"`
}

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
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
			{Name: "country_id", Type: proto.ColumnType_INT, Description: "The ID of the country associated with the organization."},
			{Name: "timezone_id", Type: proto.ColumnType_INT, Description: "The ID of the timezone associated with the organization."},
			{Name: "license", Type: proto.ColumnType_JSON, Description: "Licence information and limits."},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: "Zone where the organization exists."},
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

	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var id = int(d.KeyColumnQuals["id"].GetInt64Value())
	var config = client.NewRequestConfig("organizations", id)
	utils.ColumnsToParams(&config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})

	var result = &OrganizationResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		logger.Info("getOrganization", err.Error())
		return nil, err
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

	var config = client.NewRequestConfig("organizations", 0)
	utils.ColumnsToParams(&config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	var pagesLeft = true
	for pagesLeft {
		var result = &OrganizationListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			logger.Error("make_organization.listOrganizations", "connection_error", err)
			return nil, err
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
