package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type Registry struct {
	sync.Mutex
	gatherers  map[prometheus.Gatherer]struct{}
	flushQueue map[prometheus.Gatherer]struct{}
}

func NewRegistry() *Registry {
	return &Registry{
		gatherers:  make(map[prometheus.Gatherer]struct{}, 0),
		flushQueue: make(map[prometheus.Gatherer]struct{}, 0),
	}
}

func (r *Registry) Register(gatherer prometheus.Gatherer) {
	r.Lock()
	r.gatherers[gatherer] = struct{}{}
	r.Unlock()
}

func (r *Registry) Unregister(gatherer prometheus.Gatherer) {
	r.Lock()
	r.flushQueue[gatherer] = struct{}{}
	delete(r.gatherers, gatherer)
	r.Unlock()
}

func (r *Registry) Gather() ([]*dto.MetricFamily, error) {
	r.Lock()

	gatherers := make(prometheus.Gatherers, len(r.gatherers)+len(r.flushQueue))

	i := 0

	for gatherer := range r.gatherers {
		gatherers[i] = gatherer
		i++
	}

	for gatherer := range r.flushQueue {
		gatherers[i] = gatherer
		i++

		delete(r.flushQueue, gatherer)
	}

	r.Unlock()

	return gatherers.Gather()
}
