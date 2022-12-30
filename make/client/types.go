package client

import "net/url"

func NewRequestConfig(endpoint string, recordId int) RequestConfig {
	return RequestConfig{
		Endpoint:   endpoint,
		RecordId:   recordId,
		Params:     url.Values{},
		Pagination: RequestPagination{},
	}
}

type RequestConfig struct {
	Endpoint   string
	RecordId   int
	Params     url.Values
	Pagination RequestPagination
}

type RequestPagination struct {
	Limit  int
	Offset int
}
