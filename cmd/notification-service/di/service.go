package di

import (
	"context"

	"github.com/boreq/errors"
	"github.com/hashicorp/go-multierror"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/ports/http"
)

type Service struct {
	app        app.Application
	server     http.Server
	downloader *app.Downloader
}

func NewService(
	app app.Application,
	server http.Server,
	downloader *app.Downloader,
) Service {
	return Service{
		app:        app,
		server:     server,
		downloader: downloader,
	}
}

func (s Service) App() app.Application {
	return s.app
}

func (s Service) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error)
	runners := 0

	runners++
	go func() {
		errCh <- s.server.ListenAndServe(ctx)
	}()

	runners++
	go func() {
		errCh <- s.downloader.Run(ctx)
	}()

	var err error
	for i := 0; i < runners; i++ {
		err = multierror.Append(err, errors.Wrap(<-errCh, "error returned by runner"))
		cancel()
	}

	return err
}
