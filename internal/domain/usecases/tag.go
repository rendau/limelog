package usecases

import (
	"context"
)

func (u *St) TagSet(ctx context.Context,
	value string) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	return u.cr.Tag.Set(ctx, value)
}

func (u *St) TagList(ctx context.Context) ([]string, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Tag.List(ctx)
}

func (u *St) TagRemove(ctx context.Context,
	value string) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	return u.cr.Tag.Remove(ctx, value)
}
