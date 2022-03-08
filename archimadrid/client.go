package archimadrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ReneKroon/ttlcache/v2"
)

const (
	DefaultURL = "https://www.archimadrid.org/index.php/oracion-y-liturgia/index.php?option=com_archimadrid&format=ajax&task=leer_lecturas"
	DefaultTTL = 24 * time.Hour
)

const ResponsePrefix = "response "

type Archimadrid interface {
	GetGospel(context.Context, time.Time) (*Gospel, error)
	GetFirstLecture(context.Context, time.Time) (*Gospel, error)
	GetSecondLecture(context.Context, time.Time) (*Gospel, error)
	GetPsalm(context.Context, time.Time) (*Gospel, error)
}

type Client struct {
	url string
	ttl time.Duration
	ttlcache.SimpleCache
}

// gospelResponse is a struct that contains the response
// from the Archimadrid API
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

// Magnificat is a struct that groups together all the
// lectures for a particular Day
type Magnificat struct {
	Day           string  `json:"day"`
	FirstLecture  *Gospel `json:"first_lecture"`
	Psalm         *Gospel `json:"psalm"`
	SecondLecture *Gospel `json:"second_lecture,omitempty"`
	Gosp          *Gospel `json:"gospel"`
}

// Option is a function to apply settings to Client structure
type Option func(c *Client) Option

// NewClient returns a new instance of Client
func NewClient(opts ...Option) *Client {
	m := &Client{
		url: DefaultURL,
		ttl: DefaultTTL,
	}

	for _, opt := range opts {
		opt(m)
	}

	cache := ttlcache.NewCache()
	cache.SetTTL(m.ttl)
	m.SimpleCache = cache

	return m
}

// SetURL Sets the URL for Client
func SetURL(url string) Option {
	return func(c *Client) Option {
		prev := c.url
		c.url = url
		return SetURL(prev)
	}
}

// SetCacheTTL Sets the TTL for the Cache
func SetCacheTTL(ttl time.Duration) Option {
	return func(c *Client) Option {
		prev := c.ttl
		c.ttl = ttl
		return SetCacheTTL(prev)
	}
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

func (c *Client) getResponseFromCache(key string) (*gospelResponse, error) {
	val, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	if gospel, ok := val.(*gospelResponse); ok {
		return gospel, nil
	}
	return nil, fmt.Errorf("no valid object of type *gospelResponse found")
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

func (c *Client) getGospelOrLecture(
	ctx context.Context,
	day time.Time,
	regexString, cachePrefix string, psalm bool,
) (*Gospel, error) {
	today := day.Format("2006-01-02")
	g, err := c.getGospelFromCache(cachePrefix + today)
	if err == nil {
		return g, nil
	}

	r, err := c.getResponseFromCache(ResponsePrefix + today)
	if err == nil {
		g, err = getGospelOrLectureFromResponse(r, regexString, psalm)
		if err != nil {
			return nil, fmt.Errorf("error getting the gospel from the response: %w", err)
		}
		return g, c.saveInCache(cachePrefix+today, g)
	}

	form := url.Values{}
	form.Add("dia", today)

	req, err := http.NewRequest(http.MethodPost, c.url, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return nil, fmt.Errorf("error building the request: %w", err)
	}

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return fmt.Errorf("error performing the request: %w", err)
		}
		defer resp.Body.Close()

		gospels := []*gospelResponse{}
		if err = json.NewDecoder(resp.Body).Decode(&gospels); err != nil {
			return err
		}
		if len(gospels) <= 0 {
			return fmt.Errorf("no gospel or lecture found for day %s", today)
		}
		err = c.saveInCache(ResponsePrefix+today, gospels[0])
		if err != nil {
			return fmt.Errorf("error saving response in cache: %w", err)
		}
		g, err = getGospelOrLectureFromResponse(gospels[0], regexString, psalm)
		if err != nil {
			return fmt.Errorf("error getting the gospel from the response: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return g, c.saveInCache(cachePrefix+today, g)
}

func getGospelOrLectureFromResponse(response *gospelResponse, regexString string, psalm bool) (*Gospel, error) {
	r := regexp.MustCompile(regexString)
	text := strings.ReplaceAll(response.PostContent, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = r.FindString(text)
	if text == "" {
		return &Gospel{Day: response.PostTitle}, nil
	}
	var replaceString string
	if replaces := r.FindStringSubmatch(text); len(replaces) > 1 {
		replaceString = replaces[1]
	} else {
		return nil, errors.New("invalid regex")
	}
	text = strings.ReplaceAll(text, replaceString, "")
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
			if !psalm && i == gospelNodes.Length()-1 {
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
