package observeutil

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProMetrics struct {
	hostname string
	HTTP     *HTTPMetrics
}

type Registry = prometheus.Registry

func Setup() ProMetrics {
	metrics := initializeMetrics()

	prometheus.MustRegister(
		metrics.HTTP.panicsTotal,
		metrics.HTTP.requestsTotal,
		metrics.HTTP.requests2xxTotal,
		metrics.HTTP.requests4xxTotal,
		metrics.HTTP.requests5xxTotal,
		metrics.HTTP.inFlightRequests,
		metrics.HTTP.requestDurationMilli,
	)

	return metrics
}

func initializeMetrics() ProMetrics {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "sheyplug-host"
	}

	return ProMetrics{
		HTTP: NewHTTPMetrics(hostname),
	}
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
