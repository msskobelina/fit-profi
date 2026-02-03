package prometheus

import (
	"github.com/msskobelina/fit-profi/pkg/metric"
	"github.com/prometheus/client_golang/prometheus"
)

type CounterVec struct {
	vec *prometheus.CounterVec
}

func NewCounterVec(
	adapter *Adapter,
	entity metric.MetricWithLabels,
) *CounterVec {
	cv := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: entity.GetName(),
			Help: entity.GetDescription(),
		},
		entity.GetLabels(),
	)

	adapter.Register(cv)

	return &CounterVec{vec: cv}
}

func (c *CounterVec) AddWithLabelValues(lvs []string, value float64) {
	c.vec.WithLabelValues(lvs...).Add(value)
}

func (c *CounterVec) IncWithNamedValues(lvs map[string]string) {
	c.vec.With(lvs).Inc()
}
