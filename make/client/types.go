package client

import "net/url"

func NewRequestConfig(endpoint string, recordId int64) RequestConfig {
	return RequestConfig{
		Endpoint:   endpoint,
		RecordId:   recordId,
		Params:     url.Values{},
		Pagination: RequestPagination{},
	}
}

type RequestConfig struct {
	Endpoint   string
	RecordId   int64
	Params     url.Values
	Pagination RequestPagination
}

type RequestPagination struct {
	Limit  int64
	Offset int64
}
