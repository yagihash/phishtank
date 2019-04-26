/*
 Package phishtank provides API access on PhishTank.
*/
package phishtank

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// URL for API call of PhishTank
	APIURL = "http://checkurl.phishtank.com/checkurl/"
	// PhishTank API is supporting XML, PHP serialized obj, and json
	// Package phishtank is only supporting JSON
	APIFORMAT = "json"
)

// httpClient defines minimal interface for client
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Param is parameter for API call
type Param struct {
	name  string
	value string
}

// ReponseMetadata is struct for metadata in API reponse
type ResponseMetadata struct {
	Timestamp string `json:"timestamp"`
	ServerID  string `json:"serverid"`
	Status    string `json:"status"`
	RequestID string `json:"requestid"`
}

// InvalidContentType is type of error when content-type of response is not JSON
type InvalidContentType struct{}

func (e *InvalidContentType) Error() string {
	return "Response is not JSON"
}

// Client is HTTP client for API call
type Client struct {
	apikey     string
	endpoint   string
	httpclient httpClient
}

// Option is client option
type Option func(*Client)

// OptionHTTPClient is option func for replacing HTTP client
func OptionHttpClient(client httpClient) func(*Client) {
	return func(c *Client) {
		c.httpclient = client
	}
}

// OptionAPIURL is option func for replacing endpoint url
func OptionAPIURL(u string) func(*Client) {
	return func(c *Client) {
		c.endpoint = u
	}
}

// New creates a phishtank client with given apikey and options
func New(apikey string, options ...Option) *Client {
	c := &Client{
		apikey:     apikey,
		endpoint:   APIURL,
		httpclient: &http.Client{},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// post sends HTTP POST request with given parameters
func (c *Client) post(params ...*Param) ([]byte, error) {
	values := url.Values{}
	values.Set("api_key", c.apikey)
	values.Set("format", APIFORMAT)

	for _, param := range params {
		values.Set(param.name, param.value)
	}

	req, err := http.NewRequest(
		"POST",
		c.endpoint,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Header.Get("content-type") != "application/json" {
		return nil, &InvalidContentType{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}
