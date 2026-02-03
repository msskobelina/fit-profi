package prometheus

import (
	"github.com/msskobelina/fit-profi/pkg/metric"
)

type MetricAdapter struct {
	adapter *Adapter
}

func NewMetricAdapter(adapter *Adapter) *MetricAdapter {
	return &MetricAdapter{
		adapter: adapter,
	}
}

func (m *MetricAdapter) CounterVec(metricEntity metric.MetricWithLabels) metric.CounterVec {
	return NewCounterVec(m.adapter, metricEntity)
}

func (m *MetricAdapter) HistogramVec(metricEntity metric.MetricWithLabels) metric.HistogramVec {
	return NewHistogramVec(
		m.adapter,
		metricEntity,
		[]float64{5, 10, 25, 50, 100, 250, 500, 1000},
	)
}
