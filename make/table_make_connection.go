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

type Connection struct {
	Id           int                `json:"id"`
	Name         string             `json:"name"`
	AccountName  string             `json:"accountName"`
	AccountLabel string             `json:"accountLabel"`
	PackageName  string             `json:"packageName"`
	Expire       string             `json:"expire"`
	Metadata     ConnectionMetadata `json:"metadata,omitempty"`
	TeamId       int                `json:"teamId"`
	Upgradeable  bool               `json:"upgradeable"`
	Scoped       bool               `json:"scoped"`
	Scopes       []ConnectionScope  `json:"scopes,omitempty"`
	AccountType  string             `json:"accountType"`
	Editable     bool               `json:"editable"`
	Uid          string             `json:"uid"`
}

type ConnectionMetadata struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type ConnectionScope struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Account string `json:"account,omitempty"`
}

type ConnectionResponse struct {
	Connection Connection `json:"connection"`
}

type ConnectionListResponse struct {
	Connections []Connection `json:"connections"`
	Pg          Pagination   `json:"pg"`
}

func tableConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_connection",
		Description: "For most apps included in Make, it is necessary to create a connection, through which Make will communicate with the given third-party service according to the settings of a specific scenario.",
		List: &plugin.ListConfig{
			Hydrate:       listConnections,
			ParentHydrate: listOrganizations,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getConnection,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The connection ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user friendly name of the connection."},
			{Name: "account_name", Type: proto.ColumnType_STRING, Description: "The system name of the application tied to this connection."},
			{Name: "account_label", Type: proto.ColumnType_STRING, Description: "The user friendly name of the application tied to this connection."},
			{Name: "account_type", Type: proto.ColumnType_STRING, Description: "Authentication type."},
			{Name: "package_name", Type: proto.ColumnType_STRING, Description: "No idea at this point, TODO."},
			{Name: "expire", Type: proto.ColumnType_TIMESTAMP, Description: "When does the connection expire?"},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Metadata attached to the connection."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "ID of the team that owns this connection."},
			{Name: "upgradeable", Type: proto.ColumnType_BOOL, Description: "Can the connection be upgraded?"},
			{Name: "scoped", Type: proto.ColumnType_BOOL, Description: "Is the connection scoped?"},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "Security scopes of the connection.", Hydrate: getConnection},
			{Name: "editable", Type: proto.ColumnType_BOOL, Description: "Can the connection be edited?"},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID of this connection."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getConnection(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("getConnection", d, h)

	var logger = utils.GetLogger()

	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// connection id - direct query [KeyColumnQuals], or column detail request [h.Item]
	var id int
	if h.Item != nil {
		id = h.Item.(Connection).Id
	} else {
		id = int(d.KeyColumnQuals["id"].GetInt64Value())
	}
	config := client.NewRequestConfig("connections", id)

	var result = &ConnectionResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		logger.Error("make_connection.getConnection", "connection_error", err)
		return nil, err
	}

	return result.Connection, nil
}

func listConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("listConnections", d, h)

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
	var teams = h.Item.(Organization).Teams
	for _, team := range teams {
		var config = client.NewRequestConfig("connections", 0)
		utils.ColumnsToParams(&config.Params, []string{"id", "name", "accountName", "accountLabel", "packageName", "expire", "metadata", "teamId", "upgradeable", "scoped", "accountType", "editable", "uid"})
		config.Params.Set("teamId", strconv.Itoa(team.Id))
		if d.QueryContext.Limit != nil {
			config.Pagination.Limit = int(*d.QueryContext.Limit)
		}

		var pagesLeft = true
		for pagesLeft {
			var result = &ConnectionListResponse{}
			err = c.Get(&config, result)
			if err != nil {
				logger.Error("make_connection.listConnections", "connection_error", err)
				return nil, err
			}

			// stream results
			for _, i := range result.Connections {
				d.StreamListItem(ctx, i)
			}

			// break both cycles if we have enough rows
			if d.QueryStatus.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}

			// pagination
			var resultCount = len(result.Connections)
			if resultCount < config.Pagination.Limit {
				pagesLeft = false
			} else {
				config.Pagination.Offset += config.Pagination.Limit
			}
		}
	}

	return nil, nil
}
