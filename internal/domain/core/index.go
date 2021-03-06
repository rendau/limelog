package core

import (
	"sync"

	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/limelog/internal/adapters/repo"
)

type St struct {
	lg       logger.Lite
	repo     repo.Repo
	testing  bool
	authPsw  string
	sesToken string

	wg sync.WaitGroup

	Config       *Config
	Session      *Session
	Profile      *Profile
	Log          *Log
	Tag          *Tag
	Notification *Notification
	Jobs         *Jobs

	nfProviders []*NotificationProviderSt
}

func New(
	lg logger.Lite,
	repo repo.Repo,
	testing bool,
	authPsw string,
	sesToken string,
) *St {
	c := &St{
		lg:       lg,
		repo:     repo,
		testing:  testing,
		authPsw:  authPsw,
		sesToken: sesToken,
	}

	c.Config = NewConfig(c)
	c.Session = NewSession(c)
	c.Profile = NewProfile(c)
	c.Log = NewLog(c)
	c.Tag = NewTag(c)
	c.Notification = NewNotification(c)
	c.Jobs = NewJobs(c)

	return c
}

func (c *St) AddProvider(v *NotificationProviderSt) {
	if len(v.Levels) == 0 {
		v.Levels = append(v.Levels, "fatal", "error", "warn")
	}

	c.nfProviders = append(c.nfProviders, v)
}

func (c *St) WaitJobs() {
	c.wg.Wait()
}
