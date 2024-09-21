package kafkaclient

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "blacklog_exporter"
)

type Metrics struct {
	Latency prometheus.HistogramVec
}

// nolint: ireturn
func newHistogramVec(histogramOpts prometheus.HistogramOpts, labelNames []string) prometheus.HistogramVec {
	ev := prometheus.NewHistogramVec(histogramOpts, labelNames)

	if err := prometheus.Register(ev); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &are); ok {
			ev, ok = are.ExistingCollector.(*prometheus.HistogramVec)
			if !ok {
				panic("different metric type registration")
			}
		} else {
			panic(err)
		}
	}

	return *ev
}

func NewMetrics() Metrics {
	latencyBuckets := []float64{
		5,
		15,
		30,
		60,
		120,  // 2 minutes
		300,  // 5 minutes
		600,  // 10 minutes
		900,  // 15 minutes
		1800, // 30 minutes
		2700, // 45 minutes
		3600, // 1 hour
		7200, // 2 hours
	}

	return Metrics{
		// nolint: exhaustruct
		Latency: newHistogramVec(prometheus.HistogramOpts{
			Namespace:   Namespace,
			Name:        "latency",
			Help:        "Latency from publish to consume duration in seconds",
			ConstLabels: nil,
			Buckets:     latencyBuckets,
		}, []string{"hostname", "output_type"}),
	}
}
