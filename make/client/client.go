package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/marekjalovec/steampipe-plugin-make/make/utils"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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
	apiKey      string
	logger      hclog.Logger
	pageSize    int
}

var clientInstance *Client

// GetClient Make client constructor
func GetClient(connection *plugin.Connection) (*Client, error) {
	if clientInstance != nil {
		return clientInstance, nil
	}

	config := getConfig(connection)
	envUrl := strings.TrimSuffix(*config.EnvironmentURL, "/")

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
		apiKey:      *config.APIKey,
		baseURL:     envUrl,
		logger:      utils.GetLogger(),
		pageSize:    defaultPageSize,
	}

	return clientInstance, nil
}

func (at *Client) rateLimit() {
	<-at.rateLimiter
}

func (at *Client) Get(config *RequestConfig, target interface{}) error {
	at.rateLimit()

	// prepare the request URL
	apiUrl := at.getApiUrl(config.Endpoint, config.RecordId)
	req, err := at.createAuthorizedRequest(apiUrl)
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

func (at *Client) getApiUrl(endpoint string, recordId int) string {
	apiUrl := fmt.Sprintf("%s/api/v2/%s", at.baseURL, endpoint)
	if recordId != 0 {
		apiUrl += fmt.Sprintf("/%d", recordId)
	}

	return apiUrl
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
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", at.apiKey))

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
	reqUrl := req.URL.RequestURI()

	// make the call
	resp, err := at.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failure on %s: %w", reqUrl, err)
	}
	defer resp.Body.Close()

	// handle HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		//return makeHTTPClientError(reqUrl, resp) // TODO
		return fmt.Errorf("HTTP request failure on %s [%d]: %w", reqUrl, resp.StatusCode, err)
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

	at.logger.Info(fmt.Sprintf("Response: %s", string(b)))

	return nil
}
