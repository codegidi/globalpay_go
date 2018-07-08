package globalpay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// APIError includes the response from the globalpay API and some HTTP request info
type APIError struct {
	Message        string        `json:"message,omitempty"`
	HTTPStatusCode int           `json:"code,omitempty"`
	Details        ErrorResponse `json:"details,omitempty"`
	URL            *url.URL      `json:"url,omitempty"`
	Header         http.Header   `json:"header,omitempty"`
}

// APIError supports the error interface
func (aerr *APIError) Error() string {
	ret, _ := json.Marshal(aerr)
	return string(ret)
}

// ErrorResponse represents an error response from the globalpay API server
type ErrorResponse struct {
	Status  bool                   `json:"status,omitempty"`
	Message string                 `json:"message,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func newAPIError(resp *http.Response) *APIError {
	p, _ := ioutil.ReadAll(resp.Body)

	var globalpayErrorResp ErrorResponse
	_ = json.Unmarshal(p, &globalpayErrorResp)
	return &APIError{
		HTTPStatusCode: resp.StatusCode,
		Header:         resp.Header,
		Details:        globalpayErrorResp,
		URL:            resp.Request.URL,
	}
}
