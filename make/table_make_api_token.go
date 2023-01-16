package make

import (
	"context"
	"github.com/marekjalovec/steampipe-plugin-make/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableApiToken(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "make_api_token",
		Description: "API tokens of the currently authenticated user.",
		List: &plugin.ListConfig{
			Hydrate: listApiTokens,
		},
		Columns: []*plugin.Column{
			// Other Columns
			{Name: "token", Type: proto.ColumnType_STRING, Description: "The user ID."},
			{Name: "label", Type: proto.ColumnType_STRING, Description: "Full name of the Role."},
			{Name: "scope", Type: proto.ColumnType_JSON, Description: "Is this Role defined in an Organization, or is it part of the account?"},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Can ths role be used on the Organization, or Team level?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Label")},
		},
	}
}

func listApiTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listApiTokens", ctx, d, h)

	// create new Make client
	c, err := client.GetClient(ctx, d.Connection)
	if err != nil {
		return nil, err
	}

	// prepare params
	var config = client.NewRequestConfig("users/me/api-tokens")
	if d.QueryContext.Limit != nil {
		config.Pagination.Limit = int(*d.QueryContext.Limit)
	}

	// fetch data
	var result = &client.ApiTokenListResponse{}
	err = c.Get(&config, result)
	if err != nil {
		plugin.Logger(ctx).Error("make_api_token.listApiTokens", "request_error", err)
		return nil, c.HandleKnownErrors(err, "user:read")
	}

	// stream results
	for _, i := range result.ApiTokens {
		d.StreamListItem(ctx, i)
	}

	return nil, nil
}
