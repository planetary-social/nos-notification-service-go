// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"context"

	firestore2 "cloud.google.com/go/firestore"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/ports/http"
)

// Injectors from wire.go:

func BuildService(contextContext context.Context, configConfig config.Config) (Service, func(), error) {
	client, err := firestore.NewClient(contextContext, configConfig)
	if err != nil {
		return Service{}, nil, err
	}
	adaptersFactoryFn := newAdaptersFactoryFn()
	transactionProvider := firestore.NewTransactionProvider(client, adaptersFactoryFn)
	saveRegistrationHandler := app.NewSaveRegistrationHandler(transactionProvider)
	commands := app.Commands{
		SaveRegistration: saveRegistrationHandler,
	}
	getRelaysHandler := app.NewGetRelaysHandler(transactionProvider)
	getPublicKeysHandler := app.NewGetPublicKeysHandler(transactionProvider)
	queries := app.Queries{
		GetRelays:     getRelaysHandler,
		GetPublicKeys: getPublicKeysHandler,
	}
	application := app.Application{
		Commands: commands,
		Queries:  queries,
	}
	server := http.NewServer(configConfig, application)
	downloader := app.NewDownloader(transactionProvider)
	service := NewService(application, server, downloader)
	return service, func() {
	}, nil
}

func buildTransactionFirestoreAdapters(client *firestore2.Client, tx *firestore2.Transaction) (app.Adapters, error) {
	registrationRepository := firestore.NewRegistrationRepository(client, tx)
	eventRepository := firestore.NewEventRepository(client, tx)
	adapters := app.Adapters{
		Registrations: registrationRepository,
		Events:        eventRepository,
	}
	return adapters, nil
}
