package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strconv"
)

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
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Connection ID."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "ID of the Team that owns this Connection."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user friendly name of the Connection."},
			{Name: "account_name", Type: proto.ColumnType_STRING, Description: "The system name of the application tied to this Connection."},
			{Name: "account_label", Type: proto.ColumnType_STRING, Description: "The user friendly name of the application tied to this Connection."},
			{Name: "account_type", Type: proto.ColumnType_STRING, Description: "Authentication type."},
			{Name: "package_name", Type: proto.ColumnType_STRING, Description: "No idea at this point, sorry."},
			{Name: "expire", Type: proto.ColumnType_TIMESTAMP, Description: "When does the Connection expire?"},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Metadata attached to the Connection."},
			{Name: "upgradeable", Type: proto.ColumnType_BOOL, Description: "Can the Connection be upgraded?"},
			{Name: "scoped", Type: proto.ColumnType_BOOL, Description: "Is the Connection scoped?"},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "Security scopes of the Connection.", Hydrate: getConnection},
			{Name: "editable", Type: proto.ColumnType_BOOL, Description: "Can the Connection be edited?"},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID of this Connection."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getConnection(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("getConnection", d, h)

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var id int
	if h.Item != nil {
		// "scopes" column detail request
		id = h.Item.(client.Connection).Id
	} else {
		// direct query
		id = int(d.KeyColumnQuals["id"].GetInt64Value())
	}
	var config = client.NewRequestConfig(fmt.Sprintf(`connections/%d`, id))

	// fetch data
	var result = &client.ConnectionResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		logger.Error("make_connection.getConnection", "request_error", err)
		return nil, c.HandleKnownErrors(err, "connections:read")
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
	var teams = h.Item.(client.Organization).Teams
	for _, team := range teams {
		// prepare params
		var config = client.NewRequestConfig("connections")
		utils.ColumnsToParams(&config.Params, []string{"id", "name", "accountName", "accountLabel", "packageName", "expire", "metadata", "teamId", "upgradeable", "scoped", "accountType", "editable", "uid"})
		config.Params.Set("teamId", strconv.Itoa(team.Id))
		if d.QueryContext.Limit != nil {
			config.Pagination.Limit = int(*d.QueryContext.Limit)
		}

		// fetch data
		var pagesLeft = true
		for pagesLeft {
			var result = &client.ConnectionListResponse{}
			err = c.Get(&config, result)
			if err != nil {
				logger.Error("make_connection.listConnections", "request_error", err)
				return nil, c.HandleKnownErrors(err, "connections:read")
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
