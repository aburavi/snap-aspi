package http

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	//"github.com/go-kit/log/level"
	"go.opentelemetry.io/otel/trace"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RequestCount = kitprometheus.NewCounterFrom(
	stdprometheus.CounterOpts{
		Namespace: "mbank",
		Subsystem: "gateway",
		Name:      "request_count",
		Help:      "Number of requests received.",
	},
	[]string{"method"},
)

var TotalRequests = kitprometheus.NewCounterFrom(
	stdprometheus.CounterOpts{
		Namespace: "mbank",
		Subsystem: "gateway",
		Name:      "http_requests_total",
		Help:      "Number of get requests.",
	},
	[]string{"method"},
)

var ResponseStatus = kitprometheus.NewCounter(stdprometheus.NewCounterVec(
	stdprometheus.CounterOpts{
		Namespace: "mbank",
		Subsystem: "gateway",
		Name:      "response_status",
		Help:      "Status of HTTP response",
	},
	[]string{"method"}),
)

var HttpDuration = kitprometheus.NewHistogram(promauto.NewHistogramVec(
	stdprometheus.HistogramOpts{
		Namespace: "mbank",
		Subsystem: "gateway",
		Name:      "http_response_time_seconds",
		Help:      "Duration of HTTP requests.",
	},
	[]string{"method"}),
)

var RequestLatency = kitprometheus.NewSummaryFrom(
	stdprometheus.SummaryOpts{
		Namespace: "mbank",
		Subsystem: "gateway",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	},
	[]string{"method"},
)

func MakeMetricsHandler(r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer) {
	r.Path("/metrics").Handler(promhttp.Handler())
}
