package archimadrid

import (
	"github.com/igvaquero18/magnifibot/utils"
)

const DefaultURL = "https://www.archimadrid.org/index.php/oracion-y-liturgia/index.php?option=com_archimadrid&format=ajax&task=leer_lecturas"

type Client struct {
	url string
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
	return func(s *Client) Option {
		prev := s.url
		s.url = url
		return SetURL(prev)
	}
}

// SetLogger Sets the Logger for Client
func SetLogger(logger utils.Logger) Option {
	return func(s *Client) Option {
		prev := s.Logger
		s.Logger = logger
		return SetLogger(prev)
	}
}
