package client

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "mqtt_blackbox_exporter"
	Subsystem = "client"
)

// Metrics has all the client metrics.
type Metrics struct {
	ConnectionErrors prometheus.Counter
	PublishErrors    prometheus.Counter
	Pings            prometheus.Counter
	Pongs            prometheus.Counter
	PingDuration     prometheus.Histogram
}

// nolint: ireturn
func newCounter(counterOpts prometheus.CounterOpts) prometheus.Counter {
	ev := prometheus.NewCounter(counterOpts)

	if err := prometheus.Register(ev); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &ev); ok {
			ev, ok = are.ExistingCollector.(prometheus.Counter)
			if !ok {
				panic("different metric type registration")
			}
		} else {
			panic(err)
		}
	}

	return ev
}

// nolint: ireturn
func newHistogram(histogramOpts prometheus.HistogramOpts) prometheus.Histogram {
	ev := prometheus.NewHistogram(histogramOpts)

	if err := prometheus.Register(ev); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &ev); ok {
			ev, ok = are.ExistingCollector.(prometheus.Histogram)
			if !ok {
				panic("different metric type registration")
			}
		} else {
			panic(err)
		}
	}

	return ev
}

func NewMetrics() Metrics {
	return Metrics{
		ConnectionErrors: newCounter(prometheus.CounterOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "connection_errors_total",
			Help:        "total number of connection errors",
			ConstLabels: nil,
		}),
		PublishErrors: newCounter(prometheus.CounterOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "publish_errors_total",
			Help:        "total number of publish errors",
			ConstLabels: nil,
		}),
		Pings: newCounter(prometheus.CounterOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "pings_total",
			Help:        "total number of published pings",
			ConstLabels: nil,
		}),
		Pongs: newCounter(prometheus.CounterOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "pongs_totla",
			Help:        "total number of received pongs",
			ConstLabels: nil,
		}),
		PingDuration: newHistogram(prometheus.HistogramOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "ping_duration_seconds",
			Help:        "from ping to pong duration in seconds",
			ConstLabels: nil,
			Buckets:     prometheus.DefBuckets,
		}),
	}
}
