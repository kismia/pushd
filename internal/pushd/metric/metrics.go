package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "pushd"

var (
	TCPConnectedClientsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "tcp_connected_clients_total",
			Help:      "Number of opened client connections.",
		},
	)
	TCPCommandsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "tcp_commands_total",
			Help:      "Number of processed commands.",
		},
	)
)
