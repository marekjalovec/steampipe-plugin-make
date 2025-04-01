package make

import (
	"context"

	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableFunctions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_function",
		Description: "Functions in Make allow you to transform data. You use functions when mapping data from one module to another. Make offers a variety of built-in functions. On top of that, you can create your own custom functions.",
		List: &plugin.ListConfig{
			Hydrate: listFunctions,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getFunction,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Function ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Function."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the Function."},
			{Name: "code", Type: proto.ColumnType_STRING, Description: "The source code of the Function.", Hydrate: getFunction},
			{Name: "args", Type: proto.ColumnType_STRING, Description: "The Function arguments."},
			{Name: "scenarios", Type: proto.ColumnType_JSON, Description: "The IDs and names of Scenarios, where the Function is used.", Hydrate: getFunction},
			{Name: "created_by_user", Type: proto.ColumnType_JSON, Description: "Function author."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time of when the Function was updated."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time of when the Function was created."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "team_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func getFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getFunction", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var id int
	if h.Item != nil {
		id = h.Item.(makesdk.Function).Id
	} else {
		id = int(d.EqualsQuals["id"].GetInt64Value())
	}
	function, err := c.GetFunction(id)
	if err != nil {
		plugin.Logger(ctx).Error("make_functions.getFunction", "request_error", err)
		return nil, err
	}

	return function, nil
}

func listFunctions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listFunctions", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_functions.listFunctions", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				flp := c.NewFunctionListPaginator(int(d.RowsRemaining(ctx)), team.Id)
				for flp.HasMorePages() {
					functions, err := flp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_functions.listFunctions", "request_error", err)
						return nil, err
					}

					for _, i := range functions {
						i.TeamId = team.Id
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
