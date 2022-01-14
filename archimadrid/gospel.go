package archimadrid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// gospelResponse is a struct that contains the response from the API
type gospelResponse struct {
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
}

// Gospel contains the Gospel for a given day
type Gospel struct {
	Day       string `json:"day"`
	Title     string `json:"title"`
	Reference string `json:"reference"`
	Content   string `json:"content"`
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

	return getGospelFromResponse(gospels[0])
}

func getGospelFromResponse(response gospelResponse) (*Gospel, error) {
	text := strings.ReplaceAll(response.PostContent, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = regexp.MustCompile(`EVANGELIO.*`).FindString(text)
	text = strings.ReplaceAll(text, "EVANGELIO", "")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return nil, fmt.Errorf("error while creating the goquery document: %w", err)
	}
	title := doc.Find(".Tit_Lectura").First().Text()
	reference := doc.Find(".Tit_Negro_Normal").First().Text()
	gospelNodes := doc.Find("p")
	content := ""
	for i := 0; i < gospelNodes.Length(); i++ {
		node := gospelNodes.Eq(i)
		nodeContent := node.Text()
		if nodeContent != "" {
			if content == "" {
				content = nodeContent
				continue
			}
			if i == gospelNodes.Length()-1 {
				content = fmt.Sprintf("%s\n\n%s", content, nodeContent)
				continue
			}
			content = fmt.Sprintf("%s\n%s", content, nodeContent)
		}
	}
	return &Gospel{
		Day:       response.PostTitle,
		Title:     title,
		Reference: reference,
		Content:   content,
	}, nil
}
