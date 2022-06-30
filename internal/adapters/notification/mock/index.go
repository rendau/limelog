package mock

import (
	"sync"

	"github.com/rendau/dop/adapters/logger"
)

type St struct {
	lg logger.Lite

	q  []map[string]any
	mu sync.Mutex
}

func New(lg logger.Lite) *St {
	return &St{
		lg: lg,
	}
}

func (o *St) Send(msg map[string]any) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.q = append(o.q, msg)
}

func (o *St) PullAll() []map[string]any {
	o.mu.Lock()
	defer o.mu.Unlock()

	q := o.q

	o.q = make([]map[string]any, 0)

	return q
}

func (o *St) Get(key string, value string) map[string]any {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, msg := range o.q {
		if v, ok := (msg[key]).(string); ok && v == value {
			return msg
		}
	}

	return nil
}

func (o *St) Clean() {
	_ = o.PullAll()
}
