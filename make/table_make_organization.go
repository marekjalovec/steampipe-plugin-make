package make

import (
	"context"
	"github.com/marekjalovec/steampipe-plugin-make/make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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
			{Name: "id", Type: proto.ColumnType_INT, Description: "The organization ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
			{Name: "countryId", Type: proto.ColumnType_INT, Description: "The ID of the country associated with the organization."},
			{Name: "timezoneId", Type: proto.ColumnType_INT, Description: "The ID of the timezone associated with the organization. "},
			{Name: "license", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "zone", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "serviceName", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "isPaused", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "externalId", Type: proto.ColumnType_STRING, Description: "Make private instances use the externalId parameter for security reasons."},
		},
	}
}

func getOrganization(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := utils.GetLogger()

	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	logger.Info("getOrganization", "KeyColumnQuals", d.KeyColumnQuals)
	logger.Info("getOrganization", "Columns", d.QueryContext.Columns)

	id := d.KeyColumnQuals["id"].GetInt64Value()

	var result = &OrganizationResponse{}
	config := client.NewRequestConfig("organizations", id)
	utils.ColumnsToParams(&config.Params, d.QueryContext.Columns)
	err = c.Get(&config, &result)
	if err != nil {
		logger.Info("getOrganization", err.Error())
		return nil, err
	}

	return result.Organization, nil
}

func listOrganizations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	logger.Info("getOrganization", "KeyColumnQuals", d.KeyColumnQuals)
	logger.Info("getOrganization", "Columns", d.QueryContext.Columns)

	config := client.NewRequestConfig("organizations", 0)
	utils.ColumnsToParams(&config.Params, d.QueryContext.Columns)
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = *d.QueryContext.Limit
	}

	pagesLeft := true
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
		resultCount := int64(len(result.Organizations))
		if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
