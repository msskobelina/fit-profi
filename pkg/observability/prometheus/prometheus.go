package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Adapter struct {
	registry *prometheus.Registry
}

func NewAdapter() *Adapter {
	r := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = r

	return &Adapter{
		registry: r,
	}
}

func (a *Adapter) Register(c prometheus.Collector) {
	a.registry.MustRegister(c)
}

func (a *Adapter) Handler() http.Handler {
	return promhttp.HandlerFor(a.registry, promhttp.HandlerOpts{})
}
