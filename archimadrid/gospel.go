package archimadrid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// gospelResponse is a struct that contains the response from the API
type gospelResponse struct {
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
}

// Gospel contains the Gospel for a given day
type Gospel struct {
	Day      string
	Title    string
	Conteºnt string
}

func (c *Client) GetGospel(day time.Time) (*Gospel, error) {
	today := day.Format("2006-01-02")
	c.Debugw("Getting gospel for day %s", today)

	gospels := []gospelResponse{}

	httpClient := http.Client{}

	form := url.Values{}
	form.Add("dia", today)

	req, err := http.NewRequest(http.MethodPost, c.url, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return nil, fmt.Errorf("error building the request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing the request: %w", err)
	}

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&gospels); err != nil {
		return nil, err
	}

	if len(gospels) <= 0 {
		return nil, fmt.Errorf("no gospel found for day %s", today)
	}

	return getGospelFromResponse(gospels[0]), nil
}

func getGospelFromResponse(response gospelResponse) *Gospel {
	day := response.PostTitle

}
