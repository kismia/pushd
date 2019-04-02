package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type GathererPool struct {
	sync.Mutex
	gatherers  map[prometheus.Gatherer]struct{}
	flushQueue map[prometheus.Gatherer]struct{}
}

func NewGathererPool() *GathererPool {
	return &GathererPool{
		gatherers:  make(map[prometheus.Gatherer]struct{}, 0),
		flushQueue: make(map[prometheus.Gatherer]struct{}, 0),
	}
}

func (p *GathererPool) Add(gatherer prometheus.Gatherer) {
	p.Lock()
	p.gatherers[gatherer] = struct{}{}
	p.Unlock()
}

func (p *GathererPool) Flush(gatherer prometheus.Gatherer) {
	p.Lock()
	p.flushQueue[gatherer] = struct{}{}
	delete(p.gatherers, gatherer)
	p.Unlock()
}

func (p *GathererPool) Gather() ([]*dto.MetricFamily, error) {
	p.Lock()

	gatherers := make(prometheus.Gatherers, len(p.gatherers)+len(p.flushQueue))

	i := 0

	for gatherer := range p.gatherers {
		gatherers[i] = gatherer
		i++
	}

	for gatherer := range p.flushQueue {
		gatherers[i] = gatherer
		i++

		delete(p.flushQueue, gatherer)
	}

	p.Unlock()

	return gatherers.Gather()
}
