package prometheus

import (
	"github.com/msskobelina/fit-profi/pkg/metric"
	"github.com/prometheus/client_golang/prometheus"
)

type HistogramVec struct {
	vec *prometheus.HistogramVec
}

func NewHistogramVec(
	adapter *Adapter,
	entity metric.MetricWithLabels,
	buckets []float64,
) *HistogramVec {
	hv := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    entity.GetName(),
			Help:    entity.GetDescription(),
			Buckets: buckets,
		},
		entity.GetLabels(),
	)

	adapter.Register(hv)

	return &HistogramVec{vec: hv}
}

func (h *HistogramVec) AddWithNamedValues(lvs map[string]string, value float64) {
	h.vec.With(lvs).Observe(value)
}
