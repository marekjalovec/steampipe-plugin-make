package utils

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

var logger hclog.Logger

func CreateLogger(ctx context.Context) {
	logger = plugin.Logger(ctx)
}

func GetLogger() hclog.Logger {
	return logger
}
