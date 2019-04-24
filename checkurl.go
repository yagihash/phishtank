package phishtank

import (
	"encoding/base64"
	"encoding/json"
)

// Results format which is contained checkurl response
type CheckURLResults struct {
	URL        string `json:"url"`
	InDatabase bool   `json:"in_database"`
}

// Response format for checkurl response
type CheckURLResponse struct {
	Meta      ResponseMetadata `json:"meta"`
	Results   CheckURLResults  `json:"results"`
	ErrorText string           `json:"errortext"`
}

// CheckURL posts URL to phishtank and fetch check results
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
