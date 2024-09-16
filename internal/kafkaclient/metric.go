package kafkaclient

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "blacklog_exporter"
	Subsystem = "client"
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
		0.001,
		0.0015,
		0.002,
		0.0025,
		0.003,
		0.0035,
		0.004,
		0.0045,
		0.005,
		0.0055,
		0.006,
		0.0065,
		0.007,
		0.0075,
		0.008,
		0.0085,
		0.009,
		0.0095,
		0.01,
		0.015,
		0.02,
		0.025,
		0.03,
		0.045,
		0.05,
		0.065,
		0.07,
		0.08,
		0.09,
		0.1,
		0.2,
		0.3,
		0.5,
		1,
	}

	return Metrics{
		// nolint: exhaustruct
		Latency: newHistogramVec(prometheus.HistogramOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "latency",
			Help:        "from publish to consume duration in seconds",
			ConstLabels: nil,
			Buckets:     latencyBuckets,
		}, []string{"host", "output"}),
	}
}
