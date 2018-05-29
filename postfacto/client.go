package postfacto

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is the postfact API client.
type Client struct {
	url string
}

// NewClient initializes a postfacto client.
func NewClient(url string) *Client {
	return &Client{url: url}
}

// ActionItem is composed of a
type ActionItem struct {
	Description string
	Member      string
}

// Get fetches a postfacto retro object.
func (c Client) Get(team string) ([]ActionItem, error) {
	resp, err := http.Get(c.url + "/retros/" + team)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	return readResponse(resp.Body)
}

func readResponse(body io.Reader) ([]ActionItem, error) {
	var response struct {
		Retro struct {
			ActionItems []struct {
				Description string `json:"description"`
			} `json:"action_items"`
		} `json:"retro"`
	}

	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}

	actionItems := make([]ActionItem, 0, len(response.Retro.ActionItems))
	for _, actionItem := range response.Retro.ActionItems {
		split := strings.SplitAfter(actionItem.Description, ")")
		desc := split[1][1:]
		member := strings.Trim(split[0], "()")
		actionItems = append(actionItems, ActionItem{Description: desc, Member: member})
	}

	return actionItems, nil
}
