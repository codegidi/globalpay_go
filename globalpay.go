package globalpay

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	version = "0.1.0"
	defaultHTTPTimeout = 60 * time.Second
	baseURL = "https://globalpay.azurewebsites.net"
	tokenURL = "http://globalpayauthserver.azurewebsites.net/connect/token"
)

type service struct {
	client *Client
}

type Client struct {
	common service
	client *http.Client
	key string
	baseURL *url.URL
	tokenURL *url.URL

	logger Logger

	Transaction  *TransactionService

	LoggingEnabled bool
	Log            Logger
}

type Logger interface {
	Printf(format string, v ...interface{})
}

type Metadata map[string]interface{}

type Response map[string]interface{}

type RequestValues url.Values

func (v RequestValues) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, 3)
	for k, val := range v {
		m[k] = val[0]
	}
	return json.Marshal(m)
}

// ListMeta is pagination metadata for paginated responses from the Globalpay API
type ListMeta struct {
	Total     int `json:"total"`
	Skipped   int `json:"skipped"`
	PerPage   int `json:"perPage"`
	Page      int `json:"page"`
	PageCount int `json:"pageCount"`
}

// NewClient creates a new Globalpay API client with the given API key
// and HTTP client, allowing overriding of the HTTP client to use.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func NewClient(key string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	u, _ := url.Parse(baseURL)
	c := &Client{
		client:         httpClient,
		key:            key,
		baseURL:        u,
		LoggingEnabled: true,
		Log:            log.New(os.Stderr, "", log.LstdFlags),
	}

	c.common.client = c
	c.Transaction = (*TransactionService)(&c.common)

	return c
}

// Call actually does the HTTP request to Paystack API
func (c *Client) Call(method, path string, body, v interface{}) error {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	u, _ := c.baseURL.Parse(path)
	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		if c.LoggingEnabled {
			c.Log.Printf("Cannot create Paystack request: %v\n", err)
		}
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.key)

	if c.LoggingEnabled {
		c.Log.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
		c.Log.Printf("POST request data %v\n", buf)
	}

	start := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if c.LoggingEnabled {
		c.Log.Printf("Completed in %v\n", time.Since(start))
	}

	defer resp.Body.Close()
	return c.decodeResponse(resp, v)
}


// Call actually does the HTTP request to Paystack API
func (c *Client) CallFormPost(method, path string, body, v interface{}) error {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	u, _ := c.tokenURL.Parse(path)
	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		if c.LoggingEnabled {
			c.Log.Printf("Cannot create Paystack request: %v\n", err)
		}
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	//req.Header.Set("Authorization", "Bearer "+c.key)

	if c.LoggingEnabled {
		c.Log.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
		c.Log.Printf("POST request data %v\n", buf)
	}

	start := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if c.LoggingEnabled {
		c.Log.Printf("Completed in %v\n", time.Since(start))
	}

	defer resp.Body.Close()
	return c.decodeResponse(resp, v)
}


func mapstruct(data interface{}, v interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           v,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(data);
	return err
}

func (c *Client) decodeResponse(httpResp *http.Response, v interface{}) error {
	var resp Response
	respBody, err := ioutil.ReadAll(httpResp.Body)
	json.Unmarshal(respBody, &resp)

	if status, _ := resp["status"].(bool); !status || httpResp.StatusCode >= 400 {
		if c.LoggingEnabled {
			c.Log.Printf("Globalpay error: %+v", err)
		}
		return newAPIError(httpResp)
	}

	if c.LoggingEnabled {
		c.Log.Printf("Globalpay response: %v\n", resp)
	}

	if data, ok := resp["data"]; ok {
		switch t := resp["data"].(type) {
		case map[string]interface{}:
			return mapstruct(data, v)
		default:
			_ = t
			return mapstruct(resp, v)
		}
	}
	// if response data does not contain data key, map entire response to v
	return mapstruct(resp, v)
}
