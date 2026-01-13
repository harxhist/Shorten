package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"be/logger"
)

var log = logger.Logger

var (
	
	ApiRequestDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "api_request_duration",
			Help: "Duration of API requests in milliseconds",
		},
		[]string{"method", "endpoint"},
	)
	ApiErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_error_count",
			Help: "Count of API errors",
		},
		[]string{"method", "endpoint", "error_code"},
	)
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of requests received",
		},
		[]string{"method", "endpoint"},
	)
	CountOfVids = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "videos_count",
			Help: "Number of videos in request",
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(ApiRequestDuration)
	prometheus.MustRegister(ApiErrorCount)
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(CountOfVids)
}

func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Error("unable to start prometheus server ", err)
		}
	}()
}