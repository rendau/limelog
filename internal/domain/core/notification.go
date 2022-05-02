package core

import (
	"sync"
	"time"

	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/interfaces"
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
	msg      map[string]interface{}
	provider interfaces.Notification
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

		tag, ok := (msg[cns.SfTagFieldName]).(string)
		if !ok {
			tag = ""
		}

		msgText, ok := (msg[cns.MessageFieldName]).(string)
		if !ok {
			msgText = ""
		}

		key := level + "___" + tag + "___" + msgText

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

			// concurrent
			{
				c.stgMu.Lock()

				if _, ok = c.stg[nfPrv.Id+"___"+key]; !ok {
					c.stg[nfPrv.Id+"___"+key] = &notificationStgItemSt{
						msg:      msg,
						provider: nfPrv.Provider,
					}
				}

				c.stgMu.Unlock()
			}
		}
	}
}

func (c *Notification) sendRoutine() {
	const chanSize = 1000

	sendIntervalDuration := NfSendSendIntervalDuration
	if c.r.testing {
		sendIntervalDuration = time.Second
	}

	jobChan := make(chan *notificationStgItemSt, chanSize)
	doneChan := make(chan bool, chanSize)

	// start workers
	for i := 0; i < NfSendWorkerCount; i++ {
		go func() {
			for item := range jobChan {
				item.provider.Send(item.msg)
				doneChan <- true
			}
		}()
	}

	var stg map[string]*notificationStgItemSt

	for {
		stg = nil

		// concurrent
		{
			c.stgMu.Lock()

			if len(c.stg) > 0 {
				stg = c.stg
				// renew stg
				c.stg = map[string]*notificationStgItemSt{}
			}

			c.stgMu.Unlock()
		}

		for _, item := range stg {
			jobChan <- item
		}

		// wait jobs
		for range stg {
			<-doneChan
		}

		time.Sleep(sendIntervalDuration)
	}
}
