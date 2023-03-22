package make

import (
	"context"
	"github.com/marekjalovec/make-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableDataStore(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_data_store",
		Description: "Data Stores are used to store data from scenarios or for transferring data in between individual scenarios or scenario runs.",
		List: &plugin.ListConfig{
			Hydrate: listDataStores,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDataStore,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Data Store ID."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Data Store."},
			{Name: "records", Type: proto.ColumnType_INT, Description: "Number of records in the Data Store."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "The current size of the data in the Data Store."},
			{Name: "max_size", Type: proto.ColumnType_INT, Description: "The maximum size of the data that the Data Store can store."},
			{Name: "data_structure_id", Type: proto.ColumnType_INT, Description: "Data structure ID.", Hydrate: getDataStore},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "team_id", Type: proto.ColumnType_INT, Description: StandardColumnDescription("virtual")},
		},
	}
}

func getDataStore(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getDataStore", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	var id int
	if h.Item != nil {
		id = h.Item.(makesdk.DataStore).Id
	} else {
		id = int(d.EqualsQuals["id"].GetInt64Value())
	}
	dataStore, err := c.GetDataStore(id)
	if err != nil {
		plugin.Logger(ctx).Error("make_data_store.getDataStore", "request_error", err)
		return nil, err
	}

	return dataStore, nil
}

func listDataStores(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listDataStores", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	olp := c.NewOrganizationListPaginator(-1)
	for olp.HasMorePages() {
		organizations, err := olp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_data_store.listDataStores", "request_error", err)
			return nil, err
		}

		for _, organization := range organizations {
			for _, team := range organization.Teams {
				dslp := c.NewDataStoreListPaginator(int(d.RowsRemaining(ctx)), team.Id)
				for dslp.HasMorePages() {
					dataStores, err := dslp.NextPage()
					if err != nil {
						plugin.Logger(ctx).Error("make_data_store.listDataStores", "request_error", err)
						return nil, err
					}

					for _, i := range dataStores {
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
