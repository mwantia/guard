package guard

import (
	"sync"

	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metricsGuardRequestMatchCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: pluginName,
	Name:      "request_match_total",
	Help:      "Counter of guarded requests",
}, []string{"server"})

var _ sync.Once
