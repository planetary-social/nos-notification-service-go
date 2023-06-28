package main

import (
	"context"
	"fmt"
	configadapters "github.com/planetary-social/go-notification-service/service/adapters/config"
	"os"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/cmd/notification-service/di"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := configadapters.NewEnvironmentConfigLoader().Load()
	if err != nil {
		return errors.Wrap(err, "error creating a config")
	}

	service, cleanup, err := di.BuildService(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "error building a service")
	}
	defer cleanup()

	return service.Run(ctx)
}
