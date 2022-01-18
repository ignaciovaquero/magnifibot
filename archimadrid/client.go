package archimadrid

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/igvaquero18/magnifibot/utils"
)

const DefaultURL = "https://www.archimadrid.org/index.php/oracion-y-liturgia/index.php?option=com_archimadrid&format=ajax&task=leer_lecturas"

type Client struct {
	url string
	utils.Logger
	ttlcache.SimpleCache
}

// Option is a function to apply settings to Client structure
type Option func(c *Client) Option

// NewClient returns a new instance of Client
func NewClient(opts ...Option) *Client {
	cache := ttlcache.NewCache()
	cache.SetTTL(24 * time.Hour) // TODO: Parameterize
	m := &Client{
		url:         DefaultURL,
		Logger:      &utils.DefaultLogger{},
		SimpleCache: cache,
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
