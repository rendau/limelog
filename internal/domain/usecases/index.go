package usecases

import (
	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/limelog/internal/domain/core"
)

type St struct {
	lg logger.Lite

	cr *core.St
}

func New(
	lg logger.Lite,
	cr *core.St,
) *St {
	return &St{
		lg: lg,
		cr: cr,
	}
}
