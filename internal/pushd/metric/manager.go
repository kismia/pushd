package metric

import (
	"fmt"

	"github.com/kismia/pushd/internal/pushd/label"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rs/xid"
)

const clientUUIDLabel = "client_uuid"

type Manager struct {
	uuid       string
	registry   *prometheus.Registry
	collectors map[string]prometheus.Collector
}

func NewManager() *Manager {
	return &Manager{
		uuid:       xid.New().String(),
		registry:   prometheus.NewRegistry(),
		collectors: make(map[string]prometheus.Collector),
	}
}

func (m *Manager) Counter(name string, labels label.Labels) (prometheus.Counter, error) {
	labels = append(labels, label.New(clientUUIDLabel, m.uuid))

	collector, ok := m.collectors[name]
	if ok {
		counterVec, ok := collector.(*prometheus.CounterVec)
		if !ok {
			return nil, fmt.Errorf("unable to cast %#v of type %T to CounterVec", collector, collector)
		}

		return counterVec.GetMetricWithLabelValues(labels.Values()...)
	}

	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
	}, labels.Names())

	if err := m.registry.Register(counterVec); err != nil {
		return nil, err
	}

	m.collectors[name] = counterVec

	return counterVec.GetMetricWithLabelValues(labels.Values()...)
}

func (m *Manager) Gauge(name string, labels label.Labels) (prometheus.Gauge, error) {
	labels = append(labels, label.New(clientUUIDLabel, m.uuid))

	collector, ok := m.collectors[name]
	if ok {
		gaugeVec, ok := collector.(*prometheus.GaugeVec)
		if !ok {
			return nil, fmt.Errorf("unable to cast %#v of type %T to GaugeVec", collector, collector)
		}

		return gaugeVec.GetMetricWithLabelValues(labels.Values()...)
	}

	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
	}, labels.Names())

	if err := m.registry.Register(gaugeVec); err != nil {
		return nil, err
	}

	m.collectors[name] = gaugeVec

	return gaugeVec.GetMetricWithLabelValues(labels.Values()...)
}

func (m *Manager) Histogram(name string, labels label.Labels) (prometheus.Observer, error) {
	labels = append(labels, label.New(clientUUIDLabel, m.uuid))

	collector, ok := m.collectors[name]
	if ok {
		histogramVec, ok := collector.(*prometheus.HistogramVec)
		if !ok {
			return nil, fmt.Errorf("unable to cast %#v of type %T to HistogramVec", collector, collector)
		}

		return histogramVec.GetMetricWithLabelValues(labels.Values()...)
	}

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name,
	}, labels.Names())

	if err := m.registry.Register(histogramVec); err != nil {
		return nil, err
	}

	m.collectors[name] = histogramVec

	return histogramVec.GetMetricWithLabelValues(labels.Values()...)
}

func (m *Manager) Summary(name string, labels label.Labels) (prometheus.Observer, error) {
	labels = append(labels, label.New(clientUUIDLabel, m.uuid))

	collector, ok := m.collectors[name]
	if ok {
		summaryVec, ok := collector.(*prometheus.SummaryVec)
		if !ok {
			return nil, fmt.Errorf("unable to cast %#v of type %T to SummaryVec", collector, collector)
		}

		return summaryVec.GetMetricWithLabelValues(labels.Values()...)
	}

	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: name,
	}, labels.Names())

	if err := m.registry.Register(summaryVec); err != nil {
		return nil, err
	}

	m.collectors[name] = summaryVec

	return summaryVec.GetMetricWithLabelValues(labels.Values()...)
}

func (m *Manager) Gather() ([]*dto.MetricFamily, error) {
	return m.registry.Gather()
}
