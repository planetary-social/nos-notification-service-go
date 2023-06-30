package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	config config.Config
	logger logging.Logger
}

func NewMetricsServer(
	config config.Config,
	logger logging.Logger,
) MetricsServer {
	return MetricsServer{
		config: config,
		logger: logger.New("metricsServer"),
	}
}

func (s *MetricsServer) ListenAndServe(ctx context.Context) error {
	mux := s.createMux()

	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(ctx, "tcp", s.config.MetricsListenAddress())
	if err != nil {
		return errors.Wrap(err, "error listening")
	}

	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			fmt.Println("error closing listener:", err)
		}
	}()

	return http.Serve(listener, mux)
}

func (s *MetricsServer) createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
