package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	CountTelegramMessages = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "count_telegram_messages",
		Help: "Total number of incoming telegram messages",
	})
	CountErrorResponses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "count_error_responses",
		Help: "Total number error responses to telegram messages",
	})
)

func New() []prometheus.Collector {
	out := []prometheus.Collector{
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		CountTelegramMessages,
		CountErrorResponses,
	}
	return out
}
