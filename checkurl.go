package phishtank

import (
	"encoding/base64"
	"encoding/json"
)

type CheckURLResults struct {
	URL        string `json:"url"`
	InDatabase bool   `json:"in_database"`
}

type CheckURLResponse struct {
	Meta      ResponseMetadata `json:"meta"`
	Results   CheckURLResults  `json:"results"`
	ErrorText string           `json:"errortext"`
}

func (c *Client) CheckURL(u string) (*CheckURLResponse, error) {
	param := &Param{
		name:  "url",
		value: base64.StdEncoding.EncodeToString([]byte(u)),
	}

	response := &CheckURLResponse{}

	body, err := c.post(*param)
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}
