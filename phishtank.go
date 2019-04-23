package phishtank

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	APIURL    = "http://checkurl.phishtank.com/checkurl/"
	APIFORMAT = "json"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Param struct {
	name  string
	value string
}

type ResponseMetadata struct {
	Timestamp string `json:"timestamp"`
	ServerID  string `json:"serverid"`
	Status    string `json:"status"`
	RequestID string `json:"requestid"`
}

type InvalidContentType struct{}

func (e *InvalidContentType) Error() string {
	return "Response is not JSON"
}

type Client struct {
	apikey     string
	endpoint   string
	log        *log.Logger
	httpclient httpClient
}

type Option func(*Client)

func OptionHTTPClient(client httpClient) func(*Client) {
	return func(c *Client) {
		c.httpclient = client
	}
}

func OptionLog(l log.Logger) func(*Client) {
	return func(c *Client) {
		c.log = &l
	}
}

func OptionAPIURL(u string) func(*Client) {
	return func(c *Client) {
		c.endpoint = u
	}
}

func New(apikey string, options ...Option) *Client {
	c := &Client{
		apikey:     apikey,
		endpoint:   APIURL,
		httpclient: &http.Client{},
		log:        log.New(os.Stderr, "yagihashoo/phishtank", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Client) post(params ...Param) ([]byte, error) {
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
