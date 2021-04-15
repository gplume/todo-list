package prome

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Vars exports struct
var Vars *AppVars

// AppVars contains all counters and histograms
type AppVars struct {
	ListCount             prometheus.Counter
	GetCount              prometheus.Counter
	PostCount             prometheus.Counter
	UpdateCount           prometheus.Counter
	DeleteCount           prometheus.Counter
	HTTPResponseLatencies *prometheus.HistogramVec
}

func init() {
	InitPromeVars()
	prometheus.MustRegister(Vars.ListCount)
	prometheus.MustRegister(Vars.GetCount)
	prometheus.MustRegister(Vars.PostCount)
	prometheus.MustRegister(Vars.UpdateCount)
	prometheus.MustRegister(Vars.DeleteCount)
	prometheus.MustRegister(Vars.HTTPResponseLatencies)
}

// InitPromeVars launch Prometheus vars initialization
func InitPromeVars() {
	Vars = &AppVars{
		ListCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "list_total",
			Help: "Number of full todo list successfully processed.",
		}),

		GetCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "single_get_total",
			Help: "Number of single get todo successfully processed.",
		}),

		PostCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "post_total",
			Help: "Number of added todo successfully processed.",
		}),

		UpdateCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "update_total",
			Help: "Number of updated todo successfully processed.",
		}),

		DeleteCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "delete_total",
			Help: "Number of deleted todo successfully processed.",
		}),

		HTTPResponseLatencies: prometheus.NewHistogramVec(
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
