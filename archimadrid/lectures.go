package archimadrid

import (
	"context"
	"time"
)

func (c *Client) getLectureFromCache(key string) (string, error) {
	return "", nil
}

func (c *Client) GetFirstLecture(ctx context.Context, day time.Time) (string, error) {
	today := day.Format("2006-01-02")
	return c.getLectureFromCache("first lecture " + today)
}
