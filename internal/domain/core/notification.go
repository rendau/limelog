package core

import (
	"sync"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/interfaces"
)

const (
	NfMsgChannelSize           = 10000
	NfMsgHandleWorkerCount     = 10
	NfSendSendIntervalDuration = 20 * time.Second
	NfSendWorkerCount          = 10
)

type Notification struct {
	r *St

	nfMsgCh chan map[string]interface{}

	stg   map[string]*notificationStgItemSt
	stgMu sync.Mutex
}

type notificationStgItemSt struct {
	createdAt time.Time
	msg       map[string]interface{}
	providers map[string]interfaces.Notification
}

func NewNotification(r *St) *Notification {
	res := &Notification{
		r:       r,
		nfMsgCh: make(chan map[string]interface{}, NfMsgChannelSize),
		stg:     map[string]*notificationStgItemSt{},
	}

	for i := 0; i < NfMsgHandleWorkerCount; i++ {
		go res.handleNotificationRoutine()
	}

	go res.sendRoutine()

	return res
}

func (c *Notification) HandleMsg(msg map[string]interface{}) {
	c.nfMsgCh <- msg
}

func (c *Notification) handleNotificationRoutine() {
	for msg := range c.nfMsgCh {
		if len(c.r.nfProviders) == 0 {
			continue
		}

		level, ok := (msg[cns.LevelFieldName]).(string)
		if !ok {
			continue
		}

		sfMsgText, ok := (msg[cns.SfMessageFieldName]).(string)
		if !ok {
			sfMsgText = ""
		}

		for _, nfPrv := range c.r.nfProviders {
			levelFound := false

			for _, lvl := range nfPrv.Levels {
				if lvl == level {
					levelFound = true
					break
				}
			}

			if !levelFound {
				continue
			}

			c.stgMu.Lock()

			if stgItem, ok := c.stg[sfMsgText]; ok {
				stgItem.providers[nfPrv.Id] = nfPrv.Provider
			} else {
				c.stg[sfMsgText] = &notificationStgItemSt{
					createdAt: time.Now(),
					msg:       msg,
					providers: map[string]interfaces.Notification{
						nfPrv.Id: nfPrv.Provider,
					},
				}
			}

			c.stgMu.Unlock()
		}
	}
}

func (c *Notification) sendRoutine() {
	type jobSt struct {
		prv interfaces.Notification
		msg map[string]interface{}
	}

	jobChan := make(chan *jobSt, 1000)
	doneChan := make(chan bool, 1000)

	// start workers
	for i := 0; i < NfSendWorkerCount; i++ {
		go func() {
			for job := range jobChan {
				job.prv.Send(job.msg)
				doneChan <- true
			}
		}()
	}

	var readyItemKeys []string
	var readyItems []*notificationStgItemSt

	for {
		readyItemKeys = nil
		readyItems = nil

		c.stgMu.Lock()

		for k, item := range c.stg {
			if time.Since(item.createdAt) > NfSendSendIntervalDuration {
				readyItemKeys = append(readyItemKeys, k)
				readyItems = append(readyItems, item)
			}
		}

		for _, k := range readyItemKeys {
			delete(c.stg, k)
		}

		c.stgMu.Unlock()

		for _, item := range readyItems {
			for _, prv := range item.providers {
				jobChan <- &jobSt{
					prv: prv,
					msg: item.msg,
				}
			}
		}

		for _, item := range readyItems {
			for range item.providers {
				<-doneChan
			}
		}

		time.Sleep(NfSendSendIntervalDuration)
	}
}
