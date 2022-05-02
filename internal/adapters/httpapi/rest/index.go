package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/rendau/limelog/internal/domain/usecases"
	"github.com/rendau/limelog/internal/interfaces"
)

type St struct {
	lg   interfaces.Logger
	ucs  *usecases.St
	cors bool

	server *http.Server
}

func New(
	lg interfaces.Logger,
	listen string,
	ucs *usecases.St,
	cors bool,
) *St {
	api := &St{
		lg:   lg,
		ucs:  ucs,
		cors: cors,
	}

	api.server = &http.Server{
		Addr:              listen,
		Handler:           api.router(),
		ReadTimeout:       2 * time.Minute,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return api
}

func (a *St) Start(eChan chan<- error) {
	go func() {
		a.lg.Infow("Start rest-api", "addr", a.server.Addr)

		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lg.Errorw("Http server closed", err)
			eChan <- err
		}
	}()
}

func (a *St) Shutdown(timeout time.Duration) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()

	return a.server.Shutdown(ctx)
}
