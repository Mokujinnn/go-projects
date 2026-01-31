package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "status"},
	)

	httpRequestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_ms",
			Help:    "HTTP request duration in ms",
			Buckets: []float64{10, 50, 100, 200, 500, 1000, 2000},
		},
	)

	httpRequestsByStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_by_status",
			Help: "HTTP requests by status",
		},
		[]string{"status"},
	)
)

func main() {
	port := 8080
	error500Prob := 0.1 // 10%
	error400Prob := 0.2 // 20%
	minTime := 10
	maxTime := 500

	mux := http.NewServeMux()

	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {

		delay := minTime + rand.Intn(maxTime-minTime)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		randValue := rand.Float64()
		var status int

		if randValue < error500Prob {
			status = 500
		} else if randValue < error500Prob+error400Prob {
			status = 400
		} else {
			status = 200
		}

		statusStr := strconv.Itoa(status)
		httpRequestsTotal.WithLabelValues(r.Method, statusStr).Inc()
		httpRequestDuration.Observe(float64(delay))
		httpRequestsByStatus.WithLabelValues(statusStr).Inc()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   status,
			"delay_ms": delay,
			"time":     time.Now().Format("15:04:05"),
		})
	})

	mux.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Server started on :%d\n", port)
	fmt.Printf("Metrics: http://localhost:%d/metrics\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
