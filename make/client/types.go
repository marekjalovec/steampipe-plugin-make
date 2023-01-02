package client

import "net/url"

func NewRequestConfig(endpoint string) RequestConfig {
	return RequestConfig{
		Endpoint:   endpoint,
		Params:     url.Values{},
		Pagination: RequestPagination{},
	}
}

type RequestConfig struct {
	Endpoint   string
	Params     url.Values
	Pagination RequestPagination
}

type RequestPagination struct {
	Limit  int
	Offset int
}
