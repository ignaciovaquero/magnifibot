package archimadrid

import (
	"context"
	"time"
)

func (c *Client) GetFirstLecture(ctx context.Context, day time.Time) (*Gospel, error) {
	return c.getGospelOrLecture(ctx, day, `(PRIMERA\sLECTURA).*?Palabra de Dios\.`, "first lecture ")
}

func (c *Client) GetSecondLecture(ctx context.Context, day time.Time) (*Gospel, error) {
	return c.getGospelOrLecture(ctx, day, `(SEGUNDA\sLECTURA).*?Palabra de Dios\.`, "second lecture ")
}
