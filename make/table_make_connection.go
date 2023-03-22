package make

import (
	"context"
	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_connection",
		Description: "For most apps included in Make, it is necessary to create a connection, through which Make will communicate with the given third-party service according to the settings of a specific scenario.",
		List: &plugin.ListConfig{
			Hydrate: listConnections,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getConnection,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Connection ID."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "ID of the Team that owns the Connection."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Connection."},
			{Name: "account_name", Type: proto.ColumnType_STRING, Description: "The system name of the application tied to the Connection."},
			{Name: "account_label", Type: proto.ColumnType_STRING, Description: "The name of the application tied to the Connection."},
			{Name: "account_type", Type: proto.ColumnType_STRING, Description: "Authentication type."},
			{Name: "package_name", Type: proto.ColumnType_STRING, Description: "Name of the Custom App to which the Connection belongs to."},
			{Name: "expire", Type: proto.ColumnType_TIMESTAMP, Description: "When does the Connection expire?"},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Metadata attached to the Connection."},
			{Name: "upgradeable", Type: proto.ColumnType_BOOL, Description: "Can the Connection be upgraded?"},
			{Name: "scoped", Type: proto.ColumnType_BOOL, Description: "Is the Connection scoped?"},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "Security scopes of the Connection.", Hydrate: getConnection, Transform: transform.FromField("Scopes").NullIfEmptySlice()},
			{Name: "editable", Type: proto.ColumnType_BOOL, Description: "Can the Connection be edited?"},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID of the Connection."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getConnection", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var id int
	if h.Item != nil {
		id = h.Item.(makesdk.Connection).Id
	} else {
		id = int(d.EqualsQuals["id"].GetInt64Value())
	}
	connection, err := c.GetConnection(id)
	if err != nil {
		plugin.Logger(ctx).Error("make_connection.getConnection", "request_error", err)
		return nil, err
	}

	return connection, nil
}

func listConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listConnections", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_connection.listConnections", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				clp := c.NewConnectionListPaginator(int(d.RowsRemaining(ctx)), team.Id)
				for clp.HasMorePages() {
					connections, err := clp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_connection.listConnections", "request_error", err)
						return nil, err
					}

					for _, i := range connections {
						d.StreamListItem(ctx, i)

						if d.RowsRemaining(ctx) <= 0 {
							return nil, nil
						}
					}
				}
			}
		}
	}

	return nil, nil
}
