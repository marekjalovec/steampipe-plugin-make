package make

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

var ConfigSchema = map[string]*schema.Attribute{
	"api_token":       {Type: schema.TypeString},
	"environment_url": {Type: schema.TypeString},
	"rate_limit":      {Type: schema.TypeInt},
}

type Config struct {
	ApiToken       *string `cty:"api_token"`
	EnvironmentUrl *string `cty:"environment_url"`
	RateLimit      *int    `cty:"rate_limit"`
}

func ConfigInstance() interface{} {
	return &Config{}
}
