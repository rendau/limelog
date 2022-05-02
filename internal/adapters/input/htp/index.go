package htp

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

func New(lg interfaces.Logger, addr string, ucs *usecases.St, cors bool) *St {
	res := &St{
		lg:   lg,
		ucs:  ucs,
		cors: cors,
	}

	res.server = &http.Server{
		Addr:              addr,
		Handler:           res.router(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return res
}

func (a *St) Start(eChan chan<- error) {
	go func() {
		a.lg.Infow("Start http-input", "addr", a.server.Addr)

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
