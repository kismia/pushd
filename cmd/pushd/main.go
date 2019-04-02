package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kismia/pushd/internal/pkg/resp"
	"github.com/kismia/pushd/internal/pushd/api"
	"github.com/kismia/pushd/internal/pushd/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tidwall/redcon"
)

type options struct {
	address        string
	metricsAddress string
	metricsPath    string
	threads        int
}

func main() {
	options := options{}
	command := &cobra.Command{
		Use:   "pushd",
		Short: "Prometheus push acceptor for ephemeral and batch jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.Run()
		},
	}

	command.Flags().StringVar(&options.address, "address", ":6379", "Gateway server address")
	command.Flags().StringVar(&options.metricsAddress, "metrics-address", ":9100", "CommandCounter server address")
	command.Flags().StringVar(&options.metricsPath, "metrics-path", "/metrics", "CommandCounter path")
	command.Flags().IntVar(&options.threads, "threads", 0, "Number of operating system threads")

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *options) Run() error {
	runtime.GOMAXPROCS(o.threads)

	gathererPool := metric.NewGathererPool()

	selfMetricsRegistry := prometheus.NewRegistry()

	selfMetricsRegistry.MustRegister(
		metric.TCPConnectedClientsTotal,
		metric.TCPCommandsTotal,
	)

	gathererPool.Add(selfMetricsRegistry)

	handler := api.NewHandler()
	respServerMux := resp.NewServeMux()
	respServerMux.HandleFunc("ping", handler.Ping)
	respServerMux.HandleFunc("quit", handler.Quit)
	respServerMux.HandleFunc("cadd", handler.CounterAdd)
	respServerMux.HandleFunc("cinc", handler.CounterInc)
	respServerMux.HandleFunc("gadd", handler.GaugeAdd)
	respServerMux.HandleFunc("gset", handler.GaugeSet)
	respServerMux.HandleFunc("gsub", handler.GaugeSub)
	respServerMux.HandleFunc("ginc", handler.GaugeInc)
	respServerMux.HandleFunc("gdec", handler.GaugeDec)
	respServerMux.HandleFunc("hist", handler.HistogramObserve)
	respServerMux.HandleFunc("summ", handler.SummaryObserve)

	respServer := redcon.NewServer(
		o.address,
		func(conn redcon.Conn, cmd redcon.Command) {
			metric.TCPCommandsTotal.Inc()

			respServerMux.ServeRESP(conn, cmd)
		},
		func(conn redcon.Conn) bool {
			metric.TCPConnectedClientsTotal.Inc()

			metricService := metric.NewService()

			gathererPool.Add(metricService)

			metric.ServiceWithContext(conn, metricService)

			return true
		},
		func(conn redcon.Conn, err error) {
			metric.TCPConnectedClientsTotal.Dec()

			gathererPool.Flush(metric.ServiceFromContext(conn))
			if err != nil {
				logrus.Error(err)
			}
		},
	)

	httpServerMux := http.NewServeMux()
	httpServerMux.Handle(o.metricsPath, promhttp.HandlerFor(gathererPool, promhttp.HandlerOpts{}))

	httpServer := &http.Server{
		Addr:    o.metricsAddress,
		Handler: httpServerMux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logrus.Println("start metrics server on", o.metricsAddress)

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		logrus.Println("start resp server on", o.address)

		if err := respServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit

	logrus.Println("received shutdown signal")

	if err := respServer.Close(); err != nil {
		logrus.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	return httpServer.Shutdown(ctx)
}
