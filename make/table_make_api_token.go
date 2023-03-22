package make

import (
	"context"
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
			{Name: "token", Type: proto.ColumnType_STRING, Description: "The API Token (masked)."},
			{Name: "label", Type: proto.ColumnType_STRING, Description: "The name of the API Token."},
			{Name: "scope", Type: proto.ColumnType_JSON, Description: "Scopes enabled for the API Token."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Creation date of the API Token."},
			{Name: "is_active", Type: proto.ColumnType_BOOL, Description: "Is the API Token currently used in make.spc?"},

			// Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: StandardColumnDescription("title"), Transform: transform.FromField("Label")},
		},
	}
}

func listApiTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listApiTokens", ctx, d, h)

	c, err := NewMakeClient(d.Connection)
	if err != nil {
		return nil, err
	}

	tlp := c.NewApiTokenListPaginator(int(d.RowsRemaining(ctx)))
	for tlp.HasMorePages() {
		tokens, err := tlp.NextPage()
		if err != nil {
			plugin.Logger(ctx).Error("make_api_token.listApiTokens", "request_error", err)
			return nil, err
		}

		// stream results
		for _, i := range tokens {
			i.IsActive = c.IsTokenActive(i.Token)
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
