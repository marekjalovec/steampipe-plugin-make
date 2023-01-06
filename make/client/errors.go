package client

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Err        error
}

func (e HttpError) Error() string {
	return fmt.Sprintf("HTTP status code %d, err: %v", e.StatusCode, e.Err)
}

func makeHttpError(url string, resp *http.Response) error {
	errorDesc := "Unknown"
	switch resp.StatusCode {
	case 404:
		errorDesc = "the requested resource wasn't found"
	}

	return &HttpError{StatusCode: resp.StatusCode, Err: fmt.Errorf(
		"HTTP request failure [%s] on %s: %s",
		resp.Status, url, errorDesc),
	}
}

func getHttpError(err error) *HttpError {
	httpErr, ok := err.(*HttpError)
	if ok {
		return httpErr
	}

	return nil
}
