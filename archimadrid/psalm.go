package archimadrid

import (
	"context"
	"time"
)

func (c *Client) GetPsalm(ctx context.Context, day time.Time) (*Gospel, error) {
	return c.getGospelOrLecture(ctx, day, `(Palabra\sde\sDios\..*<p>)<span.*?\sR.\s`, "psalm ", true)
}
