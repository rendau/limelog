package pg

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/mechta-market/limelog/internal/domain/entities"
)

func (d *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	result := &entities.ConfigSt{}

	err := d.DbQueryRow(ctx, `
		select v
		from cfg
		limit 1
	`).Scan(&result)
	if err != nil {
		if err == pgx.ErrNoRows {
			return result, nil
		}
		return nil, d.handleError(ctx, err)
	}

	return result, nil
}

func (d *St) ConfigSet(ctx context.Context, config *entities.ConfigSt) error {
	_, err := d.DbExec(ctx, `
		with u as (
			update cfg
			set v = $1
			returning 1
		)
		insert into cfg (v)
		select $1
		where not exists(select * from u)
	`, config)
	if err != nil {
		return d.handleError(ctx, err)
	}

	return nil
}
