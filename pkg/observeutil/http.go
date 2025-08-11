package observeutil

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetrics struct {
	hostname             string
	requestsTotal        *prometheus.CounterVec
	inFlightRequests     *prometheus.GaugeVec
	requests2xxTotal     *prometheus.CounterVec
	requests4xxTotal     *prometheus.CounterVec
	requests5xxTotal     *prometheus.CounterVec
	requestDurationMilli *prometheus.HistogramVec
	panicsTotal          *prometheus.CounterVec
}

func NewHTTPMetrics(hostname string) *HTTPMetrics {
	return &HTTPMetrics{

		hostname: hostname,

		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"host", "handler", "method", "status"},
		),

		requestDurationMilli: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_milliseconds",
				Help:    "Duration of HTTP requests in milliseconds",
				Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000}, // ms buckets
			},
			[]string{"host", "handler", "method"},
		),

		inFlightRequests: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests",
		}, []string{"host"}),

		requests2xxTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_2xx_total",
				Help: "Total number of 2xx client success responses",
			},
			[]string{"host", "handler", "method"},
		),

		requests4xxTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_4xx_total",
				Help: "Total number of 4xx client error responses",
			},
			[]string{"host", "handler", "method"},
		),

		requests5xxTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_5xx_total",
				Help: "Total number of 5xx server error responses",
			},
			[]string{"host", "handler", "method"},
		),

		panicsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_panics_total",
				Help: "Total number of panics recovered in HTTP handlers",
			},
			[]string{"host", "handler"},
		),
	}
}

func (m *HTTPMetrics) IncRequest(handler, method, status string) {
	m.requestsTotal.WithLabelValues(m.hostname, handler, method, status).Inc()
}

func (m *HTTPMetrics) ObserveDuration(handler, method string, durationMs float64) {
	m.requestDurationMilli.WithLabelValues(m.hostname, handler, method).Observe(durationMs)
}

func (m *HTTPMetrics) IncInFlight() {
	m.inFlightRequests.WithLabelValues(m.hostname).Inc()
}

func (m *HTTPMetrics) DecInFlight() {
	m.inFlightRequests.WithLabelValues(m.hostname).Dec()
}

func (m *HTTPMetrics) Inc2xx(handler, method string) {
	m.requests4xxTotal.WithLabelValues(m.hostname, handler, method).Inc()
}

func (m *HTTPMetrics) Inc4xx(handler, method string) {
	m.requests4xxTotal.WithLabelValues(m.hostname, handler, method).Inc()
}

func (m *HTTPMetrics) Inc5xx(handler, method string) {
	m.requests5xxTotal.WithLabelValues(m.hostname, handler, method).Inc()
}

func (m *HTTPMetrics) IncPanic(handler string) {
	m.panicsTotal.WithLabelValues(m.hostname, handler).Inc()
}
