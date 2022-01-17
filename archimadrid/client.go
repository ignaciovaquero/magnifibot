package archimadrid

import (
	"github.com/go-redis/cache/v8"
	"github.com/igvaquero18/magnifibot/utils"
)

const DefaultURL = "https://www.archimadrid.org/index.php/oracion-y-liturgia/index.php?option=com_archimadrid&format=ajax&task=leer_lecturas"

type Client struct {
	url string
	*cache.Cache
	utils.Logger
}

// Option is a function to apply settings to Client structure
type Option func(c *Client) Option

// NewClient returns a new instance of Client
func NewClient(opts ...Option) *Client {
	m := &Client{
		url:    DefaultURL,
		Logger: &utils.DefaultLogger{},
	}
	for _, opt := range opts {
		opt(m)
	}
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

// SetLogger Sets the Logger for Client
func SetLogger(logger utils.Logger) Option {
	return func(c *Client) Option {
		prev := c.Logger
		c.Logger = logger
		return SetLogger(prev)
	}
}

// SetCache sets the redis Cache in case that
// we want to enable caching
func SetCache(ca *cache.Cache) Option {
	return func(c *Client) Option {
		prev := c.Cache
		c.Cache = ca
		return SetCache(prev)
	}
}
