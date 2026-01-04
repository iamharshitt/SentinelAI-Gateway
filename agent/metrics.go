package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sentinel_requests_total",
			Help: "Total number of requests processed",
		},
	)

	actionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sentinel_actions_total",
			Help: "Total number of actions taken by the analyzer",
		},
		[]string{"action"},
	)

	requestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "sentinel_request_duration_seconds",
			Help:    "Request processing time in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, actionCounter, requestLatency)
}

// ServeMetrics starts an HTTP server exposing Prometheus metrics on addr (e.g. ":9090").
func ServeMetrics(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(addr, nil)
	return nil
}
