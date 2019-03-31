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
	flushQueueSnapshot := r.makeFlushQueueSnapshot()
	managersSnapshot := r.makeManagersSnapshot()

	metricFamilies := make([]*dto.MetricFamily, 0)

	for _, manager := range managersSnapshot {
		managerMetrics, err := manager.Gather()
		if err != nil {
			logrus.Errorln("error on gather manager metrics", err)
		} else {
			metricFamilies = appendMetrics(metricFamilies, managerMetrics)
		}
	}

	r.unregister(flushQueueSnapshot)

	return metricFamilies, nil
}

func (r *Registry) unregister(managers []*Manager) {
	r.managersMux.Lock()
	r.flushQueueMux.Lock()

	for _, manager := range managers {
		delete(r.flushQueue, manager)
		delete(r.managers, manager)
	}

	r.managersMux.Unlock()
	r.flushQueueMux.Unlock()
}

func (r *Registry) makeFlushQueueSnapshot() []*Manager {
	r.flushQueueMux.RLock()

	snapshot := make([]*Manager, len(r.flushQueue))

	i := 0

	for manager := range r.flushQueue {
		snapshot[i] = manager
		i++
	}

	r.flushQueueMux.RUnlock()

	return snapshot
}

func (r *Registry) makeManagersSnapshot() []*Manager {
	r.managersMux.RLock()

	snapshot := make([]*Manager, len(r.managers))

	i := 0

	for manager := range r.managers {
		snapshot[i] = manager
		i++
	}

	r.managersMux.RUnlock()

	return snapshot
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
