package core

import (
	"sync"

	"github.com/mechta-market/limelog/internal/interfaces"
)

type St struct {
	lg      interfaces.Logger
	db      interfaces.Db
	testing bool

	wg sync.WaitGroup

	Config *Config

	Notification *Notification
	Log          *Log
}

func New(
	lg interfaces.Logger,
	db interfaces.Db,
	testing bool,
) *St {
	c := &St{
		lg:      lg,
		db:      db,
		testing: testing,
	}

	c.Config = NewConfig(c)

	c.Notification = NewNotification(c)
	c.Log = NewLog(c)

	return c
}

func (c *St) WaitJobs() {
	c.wg.Wait()
}
