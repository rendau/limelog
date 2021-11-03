package mock

import (
	"sync"

	"github.com/mechta-market/limelog/internal/interfaces"
)

type St struct {
	lg interfaces.Logger

	q  []map[string]interface{}
	mu sync.Mutex
}

func New(lg interfaces.Logger) *St {
	return &St{
		lg: lg,
	}
}

func (o *St) Send(msg map[string]interface{}) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.q = append(o.q, msg)
}

func (o *St) PullAll() []map[string]interface{} {
	o.mu.Lock()
	defer o.mu.Unlock()

	q := o.q

	o.q = make([]map[string]interface{}, 0)

	return q
}

func (o *St) Get(key string, value string) map[string]interface{} {
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
