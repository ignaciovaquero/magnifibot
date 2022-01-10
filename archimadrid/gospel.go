package archimadrid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Gospel is a struct that contains the response from the API
type Gospel struct {
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
}

func (c *Client) GetGospel(day time.Time) (*Gospel, error) {
	httpClient := http.Client{}
	gospel := &Gospel{}
	form := url.Values{}
	form.Add("dia", time.Now().Format("2006-01-02"))
	req, err := http.NewRequest(http.MethodPost, c.url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error building the request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing the request: %w", err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(gospel); err != nil {
		return nil, err
	}
	return gospel, nil
}
