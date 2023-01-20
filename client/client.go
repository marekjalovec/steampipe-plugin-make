package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultPageSize  = 10000
	defaultRateLimit = 50
)

// Client for the Make API
type Client struct {
	client      *http.Client
	rateLimiter <-chan time.Time
	baseURL     string
	apiToken    string
	logger      hclog.Logger
	pageSize    int
	scopes      *[]string
}

var clientInstance *Client

// GetClient Make client constructor
func GetClient(ctx context.Context, connection *plugin.Connection) (*Client, error) {
	if clientInstance != nil {
		return clientInstance, nil
	}

	config, err := getConfig(connection)
	if err != nil {
		return nil, err
	}
	if config.RateLimit == nil {
		config.RateLimit = &defaultRateLimit
	}

	// rate limiter with 10% burstable rate
	var rateLimiter = make(chan time.Time, *config.RateLimit/10)
	go func() {
		for t := range time.Tick(time.Minute / time.Duration(*config.RateLimit)) {
			rateLimiter <- t
		}
	}()

	clientInstance = &Client{
		client:      http.DefaultClient,
		rateLimiter: rateLimiter,
		apiToken:    *config.ApiToken,
		baseURL:     *config.EnvironmentURL,
		logger:      plugin.Logger(ctx),
		pageSize:    defaultPageSize,
		scopes:      nil,
	}

	clientInstance.loadScopes()

	return clientInstance, nil
}

func (at *Client) rateLimit() {
	<-at.rateLimiter
}

func (at *Client) Get(config *RequestConfig, target interface{}) error {
	at.rateLimit()

	// prepare the request URL
	req, err := at.createAuthorizedRequest(fmt.Sprintf("%s/api/v2/%s", at.baseURL, config.Endpoint))
	if err != nil {
		return err
	}
	at.setQueryParams(req, config)

	// make the call
	err = at.do(req, target)
	if err != nil {
		return err
	}

	return nil
}

func (at *Client) createAuthorizedRequest(apiUrl string) (*http.Request, error) {
	at.logger.Info(fmt.Sprintf("Resource URL: %s", apiUrl))

	// make a new request
	req, err := http.NewRequestWithContext(context.Background(), "GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	// set headers and query params
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", at.apiToken))

	return req, nil
}

func (at *Client) setQueryParams(req *http.Request, config *RequestConfig) {
	// use default limit
	if config.Pagination.Limit == 0 {
		config.Pagination.Limit = at.pageSize
	}

	// set pagination params
	config.Params.Set("pg[offset]", strconv.Itoa(config.Pagination.Offset))
	config.Params.Set("pg[limit]", strconv.Itoa(config.Pagination.Limit))

	// encode params
	req.URL.RawQuery = config.Params.Encode()
	at.logger.Info(fmt.Sprintf("Query Params: %s", req.URL.RawQuery))
}

func (at *Client) do(req *http.Request, response interface{}) error {
	var reqUrl = req.URL.RequestURI()

	// make the call
	resp, err := at.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failure on %s: %w", reqUrl, err)
	}
	defer resp.Body.Close()

	// handle HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return makeHttpError(reqUrl, resp)
	}

	// read response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read error on response for %s: %w", reqUrl, err)
	}

	// parse the body
	err = json.Unmarshal(b, response)
	if err != nil {
		return fmt.Errorf("JSON decode failed on %s: %s error: %w", reqUrl, hclog.Quote(b), err)
	}

	return nil
}

func (at *Client) loadScopes() {
	var config = NewRequestConfig("users/me/api-tokens")
	var result = &ApiTokenListResponse{}
	err := at.Get(&config, result)
	if err == nil {
		for _, token := range result.ApiTokens {
			if at.IsTokenActive(token.Token) {
				at.scopes = &token.Scope
			}
		}
	}
}

func (at *Client) IsTokenActive(maskedToken string) bool {
	var parts = strings.Split(maskedToken, "-")
	return len(parts) > 0 && strings.HasPrefix(at.apiToken, parts[0])
}

func (at *Client) scopesLoaded() bool {
	return at.scopes != nil
}

func (at *Client) hasScope(scope string) bool {
	if at.scopes == nil {
		return false
	}

	for _, v := range *at.scopes {
		if v == scope {
			return true
		}
	}

	return false
}

func (at *Client) HandleKnownErrors(err error, scope string) error {
	var httpErr = getHttpError(err)
	if httpErr == nil {
		return err
	}

	// 403 Forbidden or 404 Not Found
	if httpErr.StatusCode == 403 || httpErr.StatusCode == 404 {
		return fmt.Errorf(`We couldn't fetch the resource you requested. You either don't have access to it, or it doesn't exist.`)
	}

	// 401 Unauthorized
	if httpErr.StatusCode == 401 {
		if at.scopesLoaded() && !at.hasScope(scope) {
			return fmt.Errorf(`We couldn't fetch the resource you requested, because your API Token is missing "%s" in the enabled scopes - create a new API Token and add this scope to the list, please.`, scope)
		} else {
			return fmt.Errorf(`We couldn't fetch the resource you requested. This might be caused by "%s" scope not being enabled. Check your API Token settings in Make, please.`, scope)
		}
	}

	return httpErr
}

func ToJSON(value interface{}) string {
	j, _ := json.Marshal(value)
	return string(j)
}
