package core

import (
	"sync"

	"github.com/mechta-market/limelog/internal/interfaces"
)

type St struct {
	lg       interfaces.Logger
	db       interfaces.Db
	testing  bool
	authPsw  string
	sesToken string

	wg sync.WaitGroup

	Config  *Config
	Session *Session
	Profile *Profile
	Log     *Log
}

func New(
	lg interfaces.Logger,
	db interfaces.Db,
	testing bool,
	authPsw string,
	sesToken string,
) *St {
	c := &St{
		lg:       lg,
		db:       db,
		testing:  testing,
		authPsw:  authPsw,
		sesToken: sesToken,
	}

	c.Config = NewConfig(c)
	c.Session = NewSession(c)
	c.Profile = NewProfile(c)
	c.Log = NewLog(c)

	return c
}

func (c *St) WaitJobs() {
	c.wg.Wait()
}
