package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strconv"
)

func tableDataStore(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_data_store",
		Description: "Data Stores are used to store data from scenarios or for transferring data in between individual scenarios or scenario runs.",
		List: &plugin.ListConfig{
			Hydrate:       listDataStores,
			ParentHydrate: listOrganizations,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDataStore,
		},
		Columns: []*plugin.Column{
			// Key Columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "The Data Store ID."},
			{Name: "team_id", Type: proto.ColumnType_INT, Description: "ID of the Team that owns this Data Store."},

			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user friendly name of the Data Store."},
			{Name: "records", Type: proto.ColumnType_INT, Description: "Number of records in this Data Store."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "The current size of the data that this Data Store can store."},
			{Name: "max_size", Type: proto.ColumnType_INT, Description: "The maximum size of the data that this Data Store can store."},
			{Name: "datastructure_id", Type: proto.ColumnType_INT, Description: "No idea at this point, sorry."},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Name")},
		},
	}
}

func getDataStore(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("getDataStore", ctx, d, h)

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var id = int(d.KeyColumnQuals["id"].GetInt64Value())
	var config = client.NewRequestConfig(fmt.Sprintf(`data-stores/%d`, id))
	ColumnsToParams(&config.Params, []string{"id", "name", "teamId", "records", "size", "maxSize", "datastructureId"})

	// fetch data
	var result = &client.DataStoreResponse{}
	err = c.Get(&config, &result)
	if err != nil {
		plugin.Logger(ctx).Error("make_data_store.getDataStore", "request_error", err)
		return nil, c.HandleKnownErrors(err, "datastores:read")
	}

	return result.DataStore, nil
}

func listDataStores(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listDataStores", ctx, d, h)

	if h.Item == nil {
		return nil, fmt.Errorf("parent organization not defined")
	}

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// iterate over organization teams
	var teams = h.Item.(client.Organization).Teams
	for _, team := range teams {
		// prepare params
		var config = client.NewRequestConfig("data-stores")
		ColumnsToParams(&config.Params, []string{"id", "name", "teamId", "records", "size", "maxSize", "datastructureId"})
		config.Params.Set("teamId", strconv.Itoa(team.Id))
		if d.QueryContext.Limit != nil {
			config.Pagination.Limit = int(*d.QueryContext.Limit)
		}

		// fetch data
		var pagesLeft = true
		for pagesLeft {
			var result = &client.DataStoreListResponse{}
			err = c.Get(&config, result)
			if err != nil {
				plugin.Logger(ctx).Error("make_data_store.listDataStores", "request_error", err)
				return nil, c.HandleKnownErrors(err, "datastores:read")
			}

			// stream results
			for _, i := range result.DataStores {
				d.StreamListItem(ctx, i)
			}

			// break both cycles if we have enough rows
			if d.QueryStatus.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}

			// pagination
			var resultCount = len(result.DataStores)
			if resultCount < config.Pagination.Limit {
				pagesLeft = false
			} else {
				config.Pagination.Offset += config.Pagination.Limit
			}
		}
	}

	return nil, nil
}
