package archimadrid

import (
	"context"
	"time"
)

type Lecture struct {
}

func (c *Client) getLectureFromCache(key string) (string, error) {
	val, err := c.Get(key)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func (c *Client) GetFirstLecture(ctx context.Context, day time.Time) (string, error) {
	regexString := `PRIMERA\sLECTURA.*?Palabra de Dios\.`
	today := day.Format("2006-01-02")
	return c.getLectureFromCache("first lecture " + today)
}
