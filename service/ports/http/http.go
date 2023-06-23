package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/boreq/errors"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type Server struct {
	config config.Config
	app    app.Application
	logger logging.Logger
}

func NewServer(
	config config.Config,
	app app.Application,
	logger logging.Logger,
) Server {
	return Server{
		config: config,
		app:    app,
		logger: logger.New("server"),
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	mux := s.createMux(ctx)

	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(ctx, "tcp", s.config.NostrListenAddress())
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

func (s *Server) createMux(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.serveWs(ctx, w, r)
	})
	return mux
}

func (s *Server) serveWs(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		s.logger.Error().WithError(err).Message("error upgrading the connection")
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error().WithError(err).Message("error closing the connection")
		}
	}()

	if err := s.handleConnection(ctx, conn); err != nil {
		closeErr := &websocket.CloseError{}
		if !errors.As(err, &closeErr) || closeErr.Code != websocket.CloseNormalClosure {
			s.logger.Error().WithError(err).Message("error handling the connection")
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, conn *websocket.Conn) error {
	s.logger.Debug().Message("accepted websocket connection")

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "error reading the websocket message")
		}

		message := nostr.ParseMessage(messageBytes)
		if message == nil {
			return errors.New("failed to parse the message")
		}

		switch v := message.(type) {
		case *nostr.EventEnvelope:
			event, err := domain.NewEvent(v.Event)
			if err != nil {
				return errors.Wrap(err, "error creating an event")
			}

			registration, err := domain.NewRegistrationFromEvent(event)
			if err != nil {
				return errors.Wrap(err, "error creating a registration")
			}

			cmd := app.NewSaveRegistration(
				registration,
			)

			if err := s.app.Commands.SaveRegistration.Handle(ctx, cmd); err != nil {
				return errors.Wrap(err, "error handling the registration command")
			}
		default:
			fmt.Println("received an unknown message:", message)
		}
	}
}
