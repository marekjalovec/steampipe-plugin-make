package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strconv"
)

type User struct {
	Id             int          `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Language       string       `json:"language"`
	TimezoneId     int          `json:"timezoneId"`
	LocaleId       int          `json:"localeId"`
	CountryId      int          `json:"countryId"`
	Features       UserFeatures `json:"features"`
	Avatar         string       `json:"avatar"`
	LastLogin      string       `json:"lastLogin"`
	OrganizationId int          `json:"organization_id"`
	TeamId         int          `json:"team_id"`
}

type UserFeatures struct {
	AllowApps       bool `json:"allow_apps"`
	AllowAppsJs     bool `json:"allow_apps_js"`
	PrivateModules  bool `json:"private_modules"`
	AllowAppsCommit bool `json:"allow_apps_commit"`
	LocalAccess     bool `json:"local_access"`
}

type UserListResponse struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pg"`
}

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
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Full name of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "User's email."},
			{Name: "language", Type: proto.ColumnType_STRING, Description: "Environment language."},
			{Name: "locale_id", Type: proto.ColumnType_INT, Description: "Environment locale - date and hour formats, decimal separators etc."},
			{Name: "country_id", Type: proto.ColumnType_INT, Description: "The ID of the country associated with the user."},
			{Name: "timezone_id", Type: proto.ColumnType_INT, Description: "The ID of the timezone associated with the user."},
			{Name: "features", Type: proto.ColumnType_JSON, Description: "Features enabled for the user."},
			{Name: "avatar", Type: proto.ColumnType_STRING, Description: "Gravatar URL for the user."},
			{Name: "last_login", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time of the last login."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: "Virtual column, has no data"},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "Virtual column, has no data"},
		},
	}
}

func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("listUsers", d, h)

	if h.Item == nil {
		return nil, fmt.Errorf("parent organization not defined")
	}

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// iterate over organization teams
	var organization = h.Item.(Organization)
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
			var result = &UserListResponse{}
			err = c.Get(&config, result)
			if err != nil {
				logger.Error("make_user.listUsers", "connection_error", err)
				return nil, err
			}

			// stream results
			for _, i := range result.Users {
				i.OrganizationId = organization.Id
				i.TeamId = team.Id
				d.StreamListItem(ctx, i)
			}

			// pagination
			var resultCount = len(result.Users)
			if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
				pagesLeft = false
			} else {
				config.Pagination.Offset += config.Pagination.Limit
			}
		}
	}

	return nil, nil
}
