package prometheus

import (
	"time"

	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelHandlerName          = "handlerName"
	labelRelayDownloaderState = "state"
	labelTopic                = "topic"
)

type Prometheus struct {
	applicationHandlerCallsCounter          *prometheus.CounterVec
	applicationHandlerCallDurationHistogram *prometheus.HistogramVec
	relayDownloaderStateGauge               *prometheus.GaugeVec
	subscriptionQueueLengthGauge            *prometheus.GaugeVec

	logger logging.Logger
}

func NewPrometheus(logger logging.Logger) *Prometheus {
	return &Prometheus{
		applicationHandlerCallsCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "application_handler_calls_total",
				Help: "Total number of calls to application handlers.",
			},
			[]string{labelHandlerName},
		),
		applicationHandlerCallDurationHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "application_handler_calls_duration",
				Help: "Duration of calls to application handlers in seconds.",
			},
			[]string{labelHandlerName},
		),
		relayDownloaderStateGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "relay_downloader_count",
				Help: "Number of running relay downloaders.",
			},
			[]string{labelRelayDownloaderState},
		),
		subscriptionQueueLengthGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "subscription_queue_length",
				Help: "Number of events in the subscription queue.",
			},
			[]string{labelTopic},
		),

		logger: logger.New("prometheus"),
	}
}

func (p *Prometheus) TrackApplicationCall(handlerName string) app.ApplicationCall {
	return NewApplicationCall(p, handlerName, p.logger)
}

func (p *Prometheus) MeasureRelayDownloadersState(n int, state app.RelayDownloaderState) {
	p.relayDownloaderStateGauge.With(prometheus.Labels{labelRelayDownloaderState: state.String()}).Set(float64(n))
}

func (p *Prometheus) ReportSubscriptionQueueLength(topic string, n int) {
	p.subscriptionQueueLengthGauge.With(prometheus.Labels{labelTopic: topic}).Set(float64(n))
}

type ApplicationCall struct {
	handlerName string
	p           *Prometheus
	start       time.Time
	logger      logging.Logger
}

func NewApplicationCall(p *Prometheus, handlerName string, logger logging.Logger) *ApplicationCall {
	return &ApplicationCall{
		p:           p,
		handlerName: handlerName,
		logger:      logger,
		start:       time.Now(),
	}
}

func (a *ApplicationCall) End() {
	duration := time.Since(a.start)

	a.logger.Debug().
		WithField("handlerName", a.handlerName).
		WithField("duration", duration).
		Message("application call")

	a.p.applicationHandlerCallsCounter.With(prometheus.Labels{labelHandlerName: a.handlerName}).Inc()
	a.p.applicationHandlerCallDurationHistogram.With(prometheus.Labels{labelHandlerName: a.handlerName}).Observe(duration.Seconds())
}
