package client

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type Config struct {
	APIKey         *string `cty:"api_key"`
	EnvironmentURL *string `cty:"environment_url"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"api_key":         {Type: schema.TypeString},
	"environment_url": {Type: schema.TypeString},
}

func ConfigInstance() interface{} {
	return &Config{}
}

// getConfig :: retrieve and cast connection config from query data TODO
func getConfig(connection *plugin.Connection) Config {
	if connection == nil || connection.Config == nil {
		return Config{}
	}
	config, _ := connection.Config.(Config)

	return config
}
