/*
 Package phishtank provides API access on PhishTank.
*/
package phishtank

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// URL for API call of PhishTank
	APIURL = "http://checkurl.phishtank.com/checkurl/"
	// PhishTank API is supporting XML, PHP serialized obj, and json
	// Package phishtank is only supporting JSON
	APIFORMAT = "json"

	HEADER_REQCOUNTINTERVAL = "x-request-limit-interval"
	HEADER_REQLIMIT         = "x-request-limit"
	HEADER_REQCOUNT         = "x-request-count"
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

// ResponseMetadata is struct for metadata in API response
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

/*
InvalidResponseHeader is type of error when the response lacks or contains headers related to request limit;

 X-Request-Limit-Interval
 X-Request-Limit
 X-Request-Count
*/
type InvalidResponseHeader struct {
	Name string
}

func (e *InvalidResponseHeader) Error() string {
	return "Insufficient response header"
}

// Client is HTTP client for API call
type Client struct {
	apikey               string
	endpoint             string
	httpclient           httpClient
	RequestLimitInterval uint16
	RequestLimit         uint16
	RequestCount         uint16
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

	if _, err := c.updateReqLimitInterval(resp); err != nil {
		return nil, err
	}

	if _, err := c.updateReqLimit(resp); err != nil {
		return nil, err
	}

	if _, err := c.updateReqCount(resp); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

// updateReqLimitInterval receives header value string and updates latest request limit interval of client
func (c *Client) updateReqLimitInterval(resp *http.Response) (uint16, error) {
	value := resp.Header.Get(HEADER_REQCOUNTINTERVAL)
	if value == "" {
		return 0, &InvalidResponseHeader{
			Name: HEADER_REQCOUNTINTERVAL,
		}
	}

	interval, err := strconv.ParseUint(strings.Split(value, " ")[0], 10, 16)
	c.RequestLimitInterval = uint16(interval)
	return c.RequestLimitInterval, err
}

// updateReqLimit receives header value string and updates latest request limit of client
func (c *Client) updateReqLimit(resp *http.Response) (uint16, error) {
	value := resp.Header.Get(HEADER_REQLIMIT)
	if value == "" {
		return 0, &InvalidResponseHeader{
			Name: HEADER_REQLIMIT,
		}
	}

	limit, err := strconv.ParseUint(value, 10, 16)
	c.RequestLimit = uint16(limit)
	return c.RequestLimit, err
}

// updateReqCount receives header value string and updates latest request count of client
func (c *Client) updateReqCount(resp *http.Response) (uint16, error) {
	value := resp.Header.Get(HEADER_REQCOUNT)
	if value == "" {
		return 0, &InvalidResponseHeader{
			Name: HEADER_REQCOUNT,
		}
	}

	count, err := strconv.ParseUint(value, 10, 16)
	c.RequestCount = uint16(count)
	return c.RequestCount, err
}
