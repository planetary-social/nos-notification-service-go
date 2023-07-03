// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"context"

	firestore2 "cloud.google.com/go/firestore"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/adapters/apns"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/adapters/prometheus"
	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
	"github.com/planetary-social/go-notification-service/service/ports/http"
	pubsub2 "github.com/planetary-social/go-notification-service/service/ports/pubsub"
)

// Injectors from wire.go:

func BuildService(contextContext context.Context, configConfig config.Config) (Service, func(), error) {
	logger := newLogrus()
	logrusLoggingSystem := logging.NewLogrusLoggingSystem(logger)
	loggingLogger := newSystemLogger(logrusLoggingSystem)
	client, cleanup, err := newFirestoreClient(contextContext, configConfig, loggingLogger)
	if err != nil {
		return Service{}, nil, err
	}
	adaptersFactoryFn := newAdaptersFactoryFn()
	transactionProvider := firestore.NewTransactionProvider(client, adaptersFactoryFn)
	generator := notifications.NewGenerator(loggingLogger)
	apnsAPNS, err := apns.NewAPNS(configConfig, loggingLogger)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	prometheusPrometheus := prometheus.NewPrometheus()
	processReceivedEventHandler := app.NewSaveReceivedEventHandler(transactionProvider, generator, apnsAPNS, loggingLogger, prometheusPrometheus)
	saveRegistrationHandler := app.NewSaveRegistrationHandler(transactionProvider, loggingLogger, prometheusPrometheus)
	commands := app.Commands{
		SaveReceivedEvent: processReceivedEventHandler,
		SaveRegistration:  saveRegistrationHandler,
	}
	getRelaysHandler := app.NewGetRelaysHandler(transactionProvider, prometheusPrometheus)
	getPublicKeysHandler := app.NewGetPublicKeysHandler(transactionProvider, prometheusPrometheus)
	getTokensHandler := app.NewGetTokensHandler(transactionProvider, prometheusPrometheus)
	receivedEventPubSub := pubsub.NewReceivedEventPubSub()
	getEventsHandler := app.NewGetEventsHandler(transactionProvider, receivedEventPubSub, prometheusPrometheus)
	queries := app.Queries{
		GetRelays:     getRelaysHandler,
		GetPublicKeys: getPublicKeysHandler,
		GetTokens:     getTokensHandler,
		GetEvents:     getEventsHandler,
	}
	application := app.Application{
		Commands: commands,
		Queries:  queries,
	}
	server := http.NewServer(configConfig, application, loggingLogger)
	metricsServer := http.NewMetricsServer(configConfig, loggingLogger)
	downloader := app.NewDownloader(transactionProvider, receivedEventPubSub, loggingLogger, prometheusPrometheus)
	receivedEventSubscriber := pubsub2.NewReceivedEventSubscriber(receivedEventPubSub, processReceivedEventHandler, loggingLogger)
	service := NewService(application, server, metricsServer, downloader, receivedEventSubscriber)
	return service, func() {
		cleanup()
	}, nil
}

func BuildIntegrationService(contextContext context.Context, configConfig config.Config) (Service, func(), error) {
	logger := newLogrus()
	logrusLoggingSystem := logging.NewLogrusLoggingSystem(logger)
	loggingLogger := newSystemLogger(logrusLoggingSystem)
	client, cleanup, err := newFirestoreClient(contextContext, configConfig, loggingLogger)
	if err != nil {
		return Service{}, nil, err
	}
	adaptersFactoryFn := newAdaptersFactoryFn()
	transactionProvider := firestore.NewTransactionProvider(client, adaptersFactoryFn)
	generator := notifications.NewGenerator(loggingLogger)
	apnsMock, err := apns.NewAPNSMock(configConfig, loggingLogger)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	prometheusPrometheus := prometheus.NewPrometheus()
	processReceivedEventHandler := app.NewSaveReceivedEventHandler(transactionProvider, generator, apnsMock, loggingLogger, prometheusPrometheus)
	saveRegistrationHandler := app.NewSaveRegistrationHandler(transactionProvider, loggingLogger, prometheusPrometheus)
	commands := app.Commands{
		SaveReceivedEvent: processReceivedEventHandler,
		SaveRegistration:  saveRegistrationHandler,
	}
	getRelaysHandler := app.NewGetRelaysHandler(transactionProvider, prometheusPrometheus)
	getPublicKeysHandler := app.NewGetPublicKeysHandler(transactionProvider, prometheusPrometheus)
	getTokensHandler := app.NewGetTokensHandler(transactionProvider, prometheusPrometheus)
	receivedEventPubSub := pubsub.NewReceivedEventPubSub()
	getEventsHandler := app.NewGetEventsHandler(transactionProvider, receivedEventPubSub, prometheusPrometheus)
	queries := app.Queries{
		GetRelays:     getRelaysHandler,
		GetPublicKeys: getPublicKeysHandler,
		GetTokens:     getTokensHandler,
		GetEvents:     getEventsHandler,
	}
	application := app.Application{
		Commands: commands,
		Queries:  queries,
	}
	server := http.NewServer(configConfig, application, loggingLogger)
	metricsServer := http.NewMetricsServer(configConfig, loggingLogger)
	downloader := app.NewDownloader(transactionProvider, receivedEventPubSub, loggingLogger, prometheusPrometheus)
	receivedEventSubscriber := pubsub2.NewReceivedEventSubscriber(receivedEventPubSub, processReceivedEventHandler, loggingLogger)
	service := NewService(application, server, metricsServer, downloader, receivedEventSubscriber)
	return service, func() {
		cleanup()
	}, nil
}

func buildTransactionFirestoreAdapters(client *firestore2.Client, tx *firestore2.Transaction) (app.Adapters, error) {
	relayRepository := firestore.NewRelayRepository(client, tx)
	publicKeyRepository := firestore.NewPublicKeyRepository(client, tx)
	registrationRepository := firestore.NewRegistrationRepository(client, tx, relayRepository, publicKeyRepository)
	tagRepository := firestore.NewTagRepository(client, tx)
	eventRepository := firestore.NewEventRepository(client, tx, relayRepository, tagRepository)
	adapters := app.Adapters{
		Registrations: registrationRepository,
		Events:        eventRepository,
		Relays:        relayRepository,
		PublicKeys:    publicKeyRepository,
	}
	return adapters, nil
}

// wire.go:

var downloaderSet = wire.NewSet(app.NewDownloader)

var generatorSet = wire.NewSet(notifications.NewGenerator)
