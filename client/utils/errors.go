package utils

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPError is raised when status code is not 2xx
type HTTPError struct {
	Status  int
	Message string
	Body    string
}

// This allows HTTPError to satisfy the error interface
func (he *HTTPError) Error() string {
	return fmt.Sprintf("Http Error: [%v] %v %v",
		he.Status, he.Message, he.Body)
}

// ParseHTTPStatusCodeInResponse creates a HTTPError from http status code and message
func ParseHTTPStatusCodeInResponse(response *http.Response) (*http.Response, error) {
	if response != nil && (response.StatusCode < 200 || response.StatusCode >= 400) {
		httpErr := &HTTPError{
			Status:  response.StatusCode,
			Message: response.Status,
		}
		if response.Body != nil {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return response, err
			}
			httpErr.Body = string(body)
		}
		return response, httpErr
	}
	return response, nil
}
