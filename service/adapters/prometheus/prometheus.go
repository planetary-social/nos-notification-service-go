package prometheus

import (
	"time"

	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelHandlerName = "handlerName"
)

type Prometheus struct {
	applicationHandlerCallsCounter          *prometheus.CounterVec
	applicationHandlerCallDurationHistogram *prometheus.HistogramVec
}

func NewPrometheus() *Prometheus {
	return &Prometheus{
		applicationHandlerCallsCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "application_handler_calls_total",
			},
			[]string{labelHandlerName},
		),
		applicationHandlerCallDurationHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "application_handler_calls_duration",
			},
			[]string{labelHandlerName},
		),
	}
}

func (p *Prometheus) TrackApplicationCall(handlerName string) app.ApplicationCall {
	return NewApplicationCall(p, handlerName)
}

type ApplicationCall struct {
	handlerName string
	p           *Prometheus
	start       time.Time
}

func NewApplicationCall(p *Prometheus, handlerName string) *ApplicationCall {
	return &ApplicationCall{
		p:           p,
		handlerName: handlerName,
		start:       time.Now(),
	}
}

func (a *ApplicationCall) End() {
	duration := time.Since(a.start)

	a.p.applicationHandlerCallsCounter.With(prometheus.Labels{labelHandlerName: a.handlerName}).Inc()
	a.p.applicationHandlerCallDurationHistogram.With(prometheus.Labels{labelHandlerName: a.handlerName}).Observe(duration.Seconds())
}
