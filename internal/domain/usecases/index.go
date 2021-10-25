package usecases

import (
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/mechta-market/limelog/internal/interfaces"
)

type St struct {
	lg interfaces.Logger

	db interfaces.Db
	cr *core.St
}

func New(
	lg interfaces.Logger,
	db interfaces.Db,
	cr *core.St,
) *St {
	u := &St{
		lg: lg,
		db: db,
		cr: cr,
	}

	return u
}
