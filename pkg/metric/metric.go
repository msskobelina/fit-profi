package metric

type Metric interface {
	CounterVec(metric MetricWithLabels) CounterVec
	HistogramVec(metric MetricWithLabels) HistogramVec
}

type CounterVec interface {
	AddWithLabelValues(lvs []string, value float64)
	IncWithNamedValues(lvs map[string]string)
}

type HistogramVec interface {
	AddWithNamedValues(lvs map[string]string, value float64)
}

type MetricEntity interface {
	GetName() string
	GetDescription() string
}

type MetricWithLabels interface {
	MetricEntity
	GetLabels() []string
}
