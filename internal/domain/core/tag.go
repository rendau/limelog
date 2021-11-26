package core

import (
	"context"
	"sync"
)

type Tag struct {
	r *St

	cache   map[string]bool
	cacheMu sync.Mutex
}

func NewTag(r *St) *Tag {
	return &Tag{
		r:     r,
		cache: map[string]bool{},
	}
}

func (c *Tag) Set(ctx context.Context, value string) error {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if _, ok := c.cache[value]; !ok {
		err := c.r.db.TagSet(ctx, value)
		if err != nil {
			return err
		}

		c.cache[value] = true
	}

	return nil
}

func (c *Tag) List(ctx context.Context) ([]string, error) {
	return c.r.db.TagList(ctx)
}

func (c *Tag) Remove(ctx context.Context, value string) error {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	err := c.r.db.TagRemove(ctx, value)
	if err != nil {
		return err
	}

	delete(c.cache, value)

	return nil
}

func (c *Tag) RefreshAll(ctx context.Context) error {
	finalList, err := c.r.db.LogListDistinctTag(ctx)
	if err != nil {
		return err
	}

	curList, err := c.List(ctx)
	if err != nil {
		return err
	}

	for _, fV := range finalList {
		err = c.Set(ctx, fV)
		if err != nil {
			return err
		}
	}

	for _, cV := range curList {
		found := false

		for _, fV := range finalList {
			if fV == cV {
				found = true
				break
			}
		}

		if !found {
			err = c.Remove(ctx, cV)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
