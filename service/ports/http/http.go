package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/boreq/errors"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type Server struct {
	config config.Config
	app    app.Application
}

func NewServer(config config.Config, app app.Application) Server {
	return Server{
		config: config,
		app:    app,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	mux := s.createMux()

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

func (s *Server) createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.serveWs(w, r)
	})
	return mux
}

func (s *Server) serveWs(rw http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("error upgrading the connection:", err)
		return
	}

	go func() {
		defer func() {
			err := conn.Close()
			fmt.Println("closed the connection, error:", err)
		}()

		if err := s.handleConnection(conn); err != nil {
			fmt.Println("error handling the connection:", err)
		}
	}()
}

func (s *Server) handleConnection(conn *websocket.Conn) error {
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "error reading the websocket message")
		}

		fmt.Printf("received websocket message: %s\n", string(messageBytes))

		message := nostr.ParseMessage(messageBytes)
		if message == nil {
			return errors.New("failed to parse the message")
		}

		switch v := message.(type) {
		case *nostr.EventEnvelope:
			event, err := domain.NewEventFromEnvelope(*v)
			if err != nil {
				return errors.Wrap(err, "error creating an event")
			}

			registration, err := domain.NewRegistrationFromEvent(event)
			if err != nil {
				return errors.Wrap(err, "error creating a registration")
			}

			cmd := app.NewSaveRegistration(
				event.PubKey(),
				registration,
			)

			if err := s.app.Commands.SaveRegistration.Handle(cmd); err != nil {
				return errors.Wrap(err, "error handling the registration command")
			}
		default:
			fmt.Println("received an unknown message:", message)
		}
	}
}
