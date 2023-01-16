package make

import (
	"context"
	"fmt"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableOrganizationVariable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_organization_variable",
		Description: "Organization Variables are user-set variables you can use in your scenarios.",
		List: &plugin.ListConfig{
			Hydrate:       listOrganizationVariables,
			ParentHydrate: listOrganizations,
		},
		Columns: []*plugin.Column{
			// Other Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the Organization Variable."},
			{Name: "type_id", Type: proto.ColumnType_INT, Description: "Original data type of the Organization Variable. Here, all are represented as strings in the `value` column."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of the Organization Variable, converted to string."},
			{Name: "is_system", Type: proto.ColumnType_BOOL, Description: "Is the Organization Variable set by Make, or by users?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: utils.StandardColumnDescription("title"), Transform: transform.FromField("Name")},

			// Virtual columns for the query
			{Name: "organization_id", Type: proto.ColumnType_INT, Description: utils.StandardColumnDescription("virtual")},
		},
	}
}

func listOrganizationVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	utils.LogQueryContext("listOrganizationVariables", d, h)

	if h.Item == nil {
		return nil, fmt.Errorf("no parent item found")
	}

	var logger = utils.GetLogger()

	// create new Make client
	c, err := client.GetClient(d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var orgId = h.Item.(client.Organization).Id
	var config = client.NewRequestConfig(fmt.Sprintf(`organizations/%d/variables`, orgId))
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var pagesLeft = true
	for pagesLeft {
		var result = &client.OrganizationVariableListResponse{}
		err = c.Get(&config, result)
		if err != nil {
			logger.Error("make_organization_variable.listOrganizationVariables", "request_error", err)
			return nil, c.HandleKnownErrors(err, "organization-variables:read")
		}

		// stream results
		for _, i := range result.OrganizationVariables {
			i.OrganizationId = orgId
			d.StreamListItem(ctx, i)
		}

		// pagination
		var resultCount = len(result.OrganizationVariables)
		if d.QueryStatus.RowsRemaining(ctx) <= 0 || resultCount < config.Pagination.Limit {
			pagesLeft = false
		} else {
			config.Pagination.Offset += config.Pagination.Limit
		}
	}

	return nil, nil
}
