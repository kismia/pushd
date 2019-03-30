package metric

import (
	"sync"

	dto "github.com/prometheus/client_model/go"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	managersMux   sync.RWMutex
	managers      map[*Manager]struct{}
	flushQueueMux sync.RWMutex
	flushQueue    map[*Manager]struct{}
}

func NewRegistry() *Registry {
	return &Registry{
		managers:   make(map[*Manager]struct{}, 0),
		flushQueue: make(map[*Manager]struct{}, 0),
	}
}

func (r *Registry) Register(manager *Manager) {
	r.managersMux.Lock()
	r.managers[manager] = struct{}{}
	r.managersMux.Unlock()
}

func (r *Registry) Unregister(manager *Manager) {
	r.flushQueueMux.Lock()
	r.flushQueue[manager] = struct{}{}
	r.flushQueueMux.Unlock()
}

func (r *Registry) Gather() ([]*dto.MetricFamily, error) {
	metricFamilies := make([]*dto.MetricFamily, 0)

	r.managersMux.RLock()

	for manager := range r.managers {
		managerMetrics, err := manager.Gather()
		if err != nil {
			logrus.Errorln("error on gather manager metrics", err)
		} else {
			metricFamilies = appendMetrics(metricFamilies, managerMetrics)
		}
	}

	r.managersMux.RUnlock()

	r.flushQueueMux.Lock()

	for manager := range r.flushQueue {
		delete(r.flushQueue, manager)
		delete(r.managers, manager)
	}

	r.flushQueueMux.Unlock()

	return metricFamilies, nil
}

func appendMetrics(slice, elems []*dto.MetricFamily) []*dto.MetricFamily {
	hash := make(map[string]*dto.MetricFamily)

	for i := 0; i < len(slice); i++ {
		hash[*slice[i].Name] = slice[i]
	}

	for i := 0; i < len(elems); i++ {
		if e, found := hash[*elems[i].Name]; found {
			e.Metric = append(e.Metric, elems[i].Metric...)
		} else {
			slice = append(slice, elems[i])
		}
	}

	return slice
}
