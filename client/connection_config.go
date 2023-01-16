package client

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
	"net/url"
	"strings"
)

type Config struct {
	ApiToken       *string `cty:"api_token"`
	EnvironmentURL *string `cty:"environment_url"`
	RateLimit      *int    `cty:"rate_limit"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"api_token":       {Type: schema.TypeString},
	"environment_url": {Type: schema.TypeString},
	"rate_limit":      {Type: schema.TypeInt},
}

func ConfigInstance() interface{} {
	return &Config{}
}

func getConfig(connection *plugin.Connection) (*Config, error) {
	var config Config
	var err error

	if connection == nil || connection.Config == nil {
		config = Config{}
	} else {
		config, _ = connection.Config.(Config)
	}

	// validate API token
	err = validateApiToken(config.ApiToken)
	if err != nil {
		return nil, err
	}

	// validate EnvironmentURL
	err = validateEnvironmentUrl(config.EnvironmentURL)
	if err != nil {
		return nil, err
	}
	*config.EnvironmentURL = strings.TrimSuffix(*config.EnvironmentURL, "/")

	return &config, err
}

func validateEnvironmentUrl(envUrl *string) error {
	// empty
	if envUrl == nil {
		return fmt.Errorf("[configuration - make.spc] the environment URL is not defined")
	}

	// not a valid url
	u, err := url.ParseRequestURI(*envUrl)
	if err != nil {
		return fmt.Errorf("[configuration - make.spc] the environment URL does not seem to be a properly formatted URL")
	}

	// not using https
	if strings.ToLower(u.Scheme) != "https" {
		return fmt.Errorf("[configuration - make.spc] use HTTPS protocol for the environment URL")
	}

	return nil
}

func validateApiToken(apiToken *string) error {
	// empty
	if apiToken == nil {
		return fmt.Errorf("[configuration - make.spc] the API Token is not defined; to get a token, visit the API tab in your Profile page in Make")
	}

	// invalid format
	_, err := uuid.Parse(*apiToken)
	if err != nil {
		return fmt.Errorf("[configuration - make.spc] the API Token seems to have a wrong format; to get a token, visit the API tab in your Profile page in Make")
	}

	return nil
}
