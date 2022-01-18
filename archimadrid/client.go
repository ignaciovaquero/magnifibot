package archimadrid

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/igvaquero18/magnifibot/utils"
)

const (
	DefaultURL = "https://www.archimadrid.org/index.php/oracion-y-liturgia/index.php?option=com_archimadrid&format=ajax&task=leer_lecturas"
	DefaultTTL = 24 * time.Hour
)

type Client struct {
	url string
	ttl time.Duration
	utils.Logger
	ttlcache.SimpleCache
}

// Option is a function to apply settings to Client structure
type Option func(c *Client) Option

// NewClient returns a new instance of Client
func NewClient(opts ...Option) *Client {
	m := &Client{
		url:    DefaultURL,
		Logger: &utils.DefaultLogger{},
		ttl:    DefaultTTL,
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

// SetLogger Sets the Logger for Client
func SetLogger(logger utils.Logger) Option {
	return func(c *Client) Option {
		prev := c.Logger
		c.Logger = logger
		return SetLogger(prev)
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
