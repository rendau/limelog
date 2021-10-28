package core

import (
	"context"
)

type Tag struct {
	r *St
}

func NewTag(r *St) *Tag {
	return &Tag{r: r}
}

func (c *Tag) Set(ctx context.Context, value string) error {
	return c.r.db.TagSet(ctx, value)
}

func (c *Tag) List(ctx context.Context) ([]string, error) {
	return c.r.db.TagList(ctx)
}

func (c *Tag) Remove(ctx context.Context, value string) error {
	return c.r.db.TagRemove(ctx, value)
}
