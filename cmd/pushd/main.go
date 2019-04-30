package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kismia/pushd/internal/pkg/promutil"
	"github.com/kismia/pushd/internal/pkg/resp"
	"github.com/kismia/pushd/internal/pushd/api"
	"github.com/kismia/pushd/internal/pushd/metric"
	"github.com/pkg/errors"
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
	defaultBuckets []string
	threads        int
	profiling      bool
}

func main() {
	opts := options{
		address:        ":6379",
		metricsAddress: ":9100",
		metricsPath:    "/metrics",
		defaultBuckets: promutil.BucketsToStrings(prometheus.DefBuckets),
	}

	command := &cobra.Command{
		Use:   "pushd",
		Short: "Prometheus push acceptor for ephemeral and batch jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	command.Flags().StringVar(&opts.address, "address", opts.address, "gateway server address")
	command.Flags().StringVar(&opts.metricsAddress, "metrics-address", opts.metricsAddress, "metrics server address")
	command.Flags().StringVar(&opts.metricsPath, "metrics-path", opts.metricsPath, "metrics path")
	command.Flags().StringSliceVar(&opts.defaultBuckets, "default-buckets", opts.defaultBuckets, "default histogram buckets")
	command.Flags().IntVar(&opts.threads, "threads", opts.threads, "number of operating system threads")
	command.Flags().BoolVar(&opts.profiling, "profiling", opts.profiling, "enable profiling")

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *options) Run() error {
	runtime.GOMAXPROCS(o.threads)

	defaultBuckets, err := promutil.StringsToBuckets(o.defaultBuckets)
	if err != nil {
		return errors.Wrap(err, "invalid default histogram buckets")
	}

	prometheus.DefBuckets = defaultBuckets

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

	if o.profiling {
		httpServerMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpServerMux.HandleFunc("/debug/pprof/heap", pprof.Index)
		httpServerMux.HandleFunc("/debug/pprof/mutex", pprof.Index)
		httpServerMux.HandleFunc("/debug/pprof/goroutine", pprof.Index)
		httpServerMux.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
		httpServerMux.HandleFunc("/debug/pprof/block", pprof.Index)
		httpServerMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpServerMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpServerMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpServerMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

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
