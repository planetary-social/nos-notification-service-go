package main

import (
	"context"
	"fmt"
	"os"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/cmd/notification-service/di"
	"github.com/planetary-social/go-notification-service/service/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	cfg, err := config.NewConfig(
		"",
		"test-project-id",
		"some.topic",
		os.Getenv("APNS_CERT"),
		"",
		config.EnvironmentDevelopment,
	)
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
