package archimadrid

import (
	"context"
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

func (c *Client) getGospelFromCache(key string) (*Gospel, error) {
	val, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	if gospel, ok := val.(*Gospel); ok {
		return gospel, nil
	}
	return nil, fmt.Errorf("no valid object of type *Gospel found")
}

func (c *Client) saveInCache(key string, o interface{}) error {
	return c.Set(key, o)
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c // Wait for f to return
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func (c *Client) GetGospel(ctx context.Context, day time.Time) (*Gospel, error) {
	today := day.Format("2006-01-02")
	gospel, err := c.getGospelFromCache("gospel " + today)
	if err == nil {
		return gospel, nil
	}

	form := url.Values{}
	form.Add("dia", today)

	req, err := http.NewRequest(http.MethodPost, c.url, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return nil, fmt.Errorf("error building the request: %w", err)
	}

	var g *Gospel

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return fmt.Errorf("error performing the request: %w", err)
		}
		defer resp.Body.Close()

		gospels := []gospelResponse{}
		if err = json.NewDecoder(resp.Body).Decode(&gospels); err != nil {
			return err
		}
		if len(gospels) <= 0 {
			return fmt.Errorf("no gospel found for day %s", today)
		}
		g, err = getGospelFromResponse(gospels[0])
		if err != nil {
			return fmt.Errorf("error getting the gospel from the response: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return g, c.saveInCache("gospel "+today, gospel)
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
