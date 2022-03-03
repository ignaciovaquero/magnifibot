package archimadrid

import (
	"context"
	"time"
)

func (c *Client) GetGospel(ctx context.Context, day time.Time) (*Gospel, error) {
	return c.getGospelOrLecture(ctx, day, `(EVANGELIO).*`, "gospel ", false)
}
