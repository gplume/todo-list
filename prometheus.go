package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type promeVars struct {
	listCount             prometheus.Counter
	getCount              prometheus.Counter
	postCount             prometheus.Counter
	updateCount           prometheus.Counter
	deleteCount           prometheus.Counter
	httpResponseLatencies *prometheus.HistogramVec
}

// Prometheus vars to register at startup
func registerPrometheusVars() {
	prometheus.MustRegister(app.prom.listCount)
	prometheus.MustRegister(app.prom.getCount)
	prometheus.MustRegister(app.prom.postCount)
	prometheus.MustRegister(app.prom.updateCount)
	prometheus.MustRegister(app.prom.deleteCount)
	prometheus.MustRegister(app.prom.httpResponseLatencies)
	// no err returned it just panics (Must...)
}

// Prometheus vars initialization
func initPromeVars() *promeVars {
	return &promeVars{
		listCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "list_total",
			Help: "Number of full todo list successfully processed.",
		}),

		getCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "single_get_total",
			Help: "Number of single get todo successfully processed.",
		}),

		postCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "post_total",
			Help: "Number of added todo successfully processed.",
		}),

		updateCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "update_total",
			Help: "Number of updated todo successfully processed.",
		}),

		deleteCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "delete_total",
			Help: "Number of deleted todo successfully processed.",
		}),

		httpResponseLatencies: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "todo_list_api",
				Subsystem: "http_server",
				Name:      "request_duration",
				Help:      "Distribution of http response latencies (ms), classified by code and method.",
			},
			[]string{"code", "method"},
		),
	}
}

// statsMiddleWare observe requests responses latencies on router Group (/todo) only
func statsMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		code := strconv.Itoa(c.Writer.Status())
		elapsed := time.Since(start)
		msElapsed := elapsed / time.Millisecond
		app.prom.httpResponseLatencies.WithLabelValues(code, c.Request.Method).Observe(float64(msElapsed))
	}
}
