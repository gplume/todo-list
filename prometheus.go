package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus vars to register at startup
func registerPrometheusVars() {
	prometheus.MustRegister(listCount)
	prometheus.MustRegister(getCount)
	prometheus.MustRegister(postCount)
	prometheus.MustRegister(updateCount)
	prometheus.MustRegister(deleteCount)
	prometheus.MustRegister(httpResponseLatencies)
	// no err returned it just panics (Must...)
}

// Prometheus vars initialization
var (
	listCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "list_total",
		Help: "Number of full todo list successfully processed.",
	})

	getCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "single_get_total",
		Help: "Number of single get todo successfully processed.",
	})

	postCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "post_total",
		Help: "Number of added todo successfully processed.",
	})

	updateCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_total",
		Help: "Number of updated todo successfully processed.",
	})

	deleteCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "delete_total",
		Help: "Number of deleted todo successfully processed.",
	})

	httpResponseLatencies = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "todo_list_api",
			Subsystem: "http_server",
			Name:      "request_duration",
			Help:      "Distribution of http response latencies (ms), classified by code and method.",
		},
		[]string{"code", "method"},
	)
)

// statsMiddleWare observe requests responses latencies on router Group (/todo) only
func statsMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		code := strconv.Itoa(c.Writer.Status())
		elapsed := time.Since(start)
		msElapsed := elapsed / time.Millisecond
		httpResponseLatencies.WithLabelValues(code, c.Request.Method).Observe(float64(msElapsed))
	}
}
