// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"context"

	firestore2 "cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/adapters"
	"github.com/planetary-social/go-notification-service/service/adapters/apns"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/adapters/prometheus"
	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
	"github.com/planetary-social/go-notification-service/service/ports/firestorepubsub"
	"github.com/planetary-social/go-notification-service/service/ports/http"
	"github.com/planetary-social/go-notification-service/service/ports/memorypubsub"
)

// Injectors from wire.go:

func BuildService(contextContext context.Context, configConfig config.Config) (Service, func(), error) {
	memoryEventWasAlreadySavedCache := adapters.NewMemoryEventWasAlreadySavedCache()
	logger, err := newLogger(configConfig)
	if err != nil {
		return Service{}, nil, err
	}
	client, cleanup, err := newFirestoreClient(contextContext, configConfig, logger)
	if err != nil {
		return Service{}, nil, err
	}
	watermillAdapter := logging.NewWatermillAdapter(logger)
	diBuildTransactionFirestoreAdaptersDependencies := buildTransactionFirestoreAdaptersDependencies{
		LoggerAdapter: watermillAdapter,
	}
	adaptersFactoryFn := newAdaptersFactoryFn(diBuildTransactionFirestoreAdaptersDependencies)
	transactionProvider := firestore.NewTransactionProvider(client, adaptersFactoryFn)
	prometheusPrometheus, err := prometheus.NewPrometheus(logger)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	saveReceivedEventHandler := app.NewSaveReceivedEventHandler(memoryEventWasAlreadySavedCache, transactionProvider, logger, prometheusPrometheus)
	saveRegistrationHandler := app.NewSaveRegistrationHandler(transactionProvider, logger, prometheusPrometheus)
	commands := app.Commands{
		SaveReceivedEvent: saveReceivedEventHandler,
		SaveRegistration:  saveRegistrationHandler,
	}
	getRelaysHandler := app.NewGetRelaysHandler(transactionProvider, prometheusPrometheus)
	getPublicKeysHandler := app.NewGetPublicKeysHandler(transactionProvider, prometheusPrometheus)
	getTokensHandler := app.NewGetTokensHandler(transactionProvider, prometheusPrometheus)
	receivedEventPubSub := pubsub.NewReceivedEventPubSub()
	getEventsHandler := app.NewGetEventsHandler(transactionProvider, receivedEventPubSub, prometheusPrometheus)
	getNotificationsHandler := app.NewGetNotificationsHandler(transactionProvider, prometheusPrometheus)
	queries := app.Queries{
		GetRelays:        getRelaysHandler,
		GetPublicKeys:    getPublicKeysHandler,
		GetTokens:        getTokensHandler,
		GetEvents:        getEventsHandler,
		GetNotifications: getNotificationsHandler,
	}
	application := app.Application{
		Commands: commands,
		Queries:  queries,
	}
	server := http.NewServer(configConfig, application, logger)
	metricsServer := http.NewMetricsServer(prometheusPrometheus, configConfig, logger)
	downloader := app.NewDownloader(memoryEventWasAlreadySavedCache, transactionProvider, receivedEventPubSub, logger, prometheusPrometheus)
	externalFollowChangeSubscriber, err := newExternalFollowChangeSubscriber(configConfig, watermillAdapter)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	apnsAPNS, err := apns.NewAPNS(configConfig, prometheusPrometheus, logger)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	followChangePuller := app.NewFollowChangePuller(externalFollowChangeSubscriber, apnsAPNS, queries, logger, prometheusPrometheus)
	receivedEventSubscriber := memorypubsub.NewReceivedEventSubscriber(receivedEventPubSub, saveReceivedEventHandler, logger)
	subscriber, err := firestore.NewWatermillSubscriber(client, watermillAdapter)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	generator := notifications.NewGenerator(logger)
	externalEventPublisher, err := newExternalEventPublisher(configConfig, watermillAdapter)
	if err != nil {
		cleanup()
		return Service{}, nil, err
	}
	processSavedEventHandler := app.NewProcessSavedEventHandler(transactionProvider, generator, apnsAPNS, logger, prometheusPrometheus, externalEventPublisher)
	eventSavedSubscriber := firestorepubsub.NewEventSavedSubscriber(subscriber, processSavedEventHandler, prometheusPrometheus, logger)
	service := NewService(application, server, metricsServer, downloader, followChangePuller, receivedEventSubscriber, externalFollowChangeSubscriber, eventSavedSubscriber, memoryEventWasAlreadySavedCache)
	return service, func() {
		cleanup()
	}, nil
}

func BuildIntegrationService(contextContext context.Context, configConfig config.Config) (IntegrationService, func(), error) {
	memoryEventWasAlreadySavedCache := adapters.NewMemoryEventWasAlreadySavedCache()
	logger, err := newLogger(configConfig)
	if err != nil {
		return IntegrationService{}, nil, err
	}
	client, cleanup, err := newFirestoreClient(contextContext, configConfig, logger)
	if err != nil {
		return IntegrationService{}, nil, err
	}
	watermillAdapter := logging.NewWatermillAdapter(logger)
	diBuildTransactionFirestoreAdaptersDependencies := buildTransactionFirestoreAdaptersDependencies{
		LoggerAdapter: watermillAdapter,
	}
	adaptersFactoryFn := newAdaptersFactoryFn(diBuildTransactionFirestoreAdaptersDependencies)
	transactionProvider := firestore.NewTransactionProvider(client, adaptersFactoryFn)
	prometheusPrometheus, err := prometheus.NewPrometheus(logger)
	if err != nil {
		cleanup()
		return IntegrationService{}, nil, err
	}
	saveReceivedEventHandler := app.NewSaveReceivedEventHandler(memoryEventWasAlreadySavedCache, transactionProvider, logger, prometheusPrometheus)
	saveRegistrationHandler := app.NewSaveRegistrationHandler(transactionProvider, logger, prometheusPrometheus)
	commands := app.Commands{
		SaveReceivedEvent: saveReceivedEventHandler,
		SaveRegistration:  saveRegistrationHandler,
	}
	getRelaysHandler := app.NewGetRelaysHandler(transactionProvider, prometheusPrometheus)
	getPublicKeysHandler := app.NewGetPublicKeysHandler(transactionProvider, prometheusPrometheus)
	getTokensHandler := app.NewGetTokensHandler(transactionProvider, prometheusPrometheus)
	receivedEventPubSub := pubsub.NewReceivedEventPubSub()
	getEventsHandler := app.NewGetEventsHandler(transactionProvider, receivedEventPubSub, prometheusPrometheus)
	getNotificationsHandler := app.NewGetNotificationsHandler(transactionProvider, prometheusPrometheus)
	queries := app.Queries{
		GetRelays:        getRelaysHandler,
		GetPublicKeys:    getPublicKeysHandler,
		GetTokens:        getTokensHandler,
		GetEvents:        getEventsHandler,
		GetNotifications: getNotificationsHandler,
	}
	application := app.Application{
		Commands: commands,
		Queries:  queries,
	}
	server := http.NewServer(configConfig, application, logger)
	metricsServer := http.NewMetricsServer(prometheusPrometheus, configConfig, logger)
	downloader := app.NewDownloader(memoryEventWasAlreadySavedCache, transactionProvider, receivedEventPubSub, logger, prometheusPrometheus)
	externalFollowChangeSubscriber, err := newExternalFollowChangeSubscriber(configConfig, watermillAdapter)
	if err != nil {
		cleanup()
		return IntegrationService{}, nil, err
	}
	apnsMock, err := apns.NewAPNSMock(configConfig, logger)
	if err != nil {
		cleanup()
		return IntegrationService{}, nil, err
	}
	followChangePuller := app.NewFollowChangePuller(externalFollowChangeSubscriber, apnsMock, queries, logger, prometheusPrometheus)
	receivedEventSubscriber := memorypubsub.NewReceivedEventSubscriber(receivedEventPubSub, saveReceivedEventHandler, logger)
	subscriber, err := firestore.NewWatermillSubscriber(client, watermillAdapter)
	if err != nil {
		cleanup()
		return IntegrationService{}, nil, err
	}
	generator := notifications.NewGenerator(logger)
	externalEventPublisher, err := newExternalEventPublisher(configConfig, watermillAdapter)
	if err != nil {
		cleanup()
		return IntegrationService{}, nil, err
	}
	processSavedEventHandler := app.NewProcessSavedEventHandler(transactionProvider, generator, apnsMock, logger, prometheusPrometheus, externalEventPublisher)
	eventSavedSubscriber := firestorepubsub.NewEventSavedSubscriber(subscriber, processSavedEventHandler, prometheusPrometheus, logger)
	service := NewService(application, server, metricsServer, downloader, followChangePuller, receivedEventSubscriber, externalFollowChangeSubscriber, eventSavedSubscriber, memoryEventWasAlreadySavedCache)
	integrationService := IntegrationService{
		Service:  service,
		MockAPNS: apnsMock,
	}
	return integrationService, func() {
		cleanup()
	}, nil
}

func buildTransactionFirestoreAdapters(client *firestore2.Client, tx *firestore2.Transaction, deps buildTransactionFirestoreAdaptersDependencies) (app.Adapters, error) {
	relayRepository := firestore.NewRelayRepository(client, tx)
	publicKeyRepository := firestore.NewPublicKeyRepository(client, tx)
	registrationRepository := firestore.NewRegistrationRepository(client, tx, relayRepository, publicKeyRepository)
	tagRepository := firestore.NewTagRepository(client, tx)
	eventRepository := firestore.NewEventRepository(client, tx, relayRepository, tagRepository)
	loggerAdapter := deps.LoggerAdapter
	publisher, err := firestore.NewWatermillPublisher(client, loggerAdapter)
	if err != nil {
		return app.Adapters{}, err
	}
	firestorePublisher := firestore.NewPublisher(publisher, tx)
	appAdapters := app.Adapters{
		Registrations: registrationRepository,
		Relays:        relayRepository,
		PublicKeys:    publicKeyRepository,
		Events:        eventRepository,
		Tags:          tagRepository,
		Publisher:     firestorePublisher,
	}
	return appAdapters, nil
}

// wire.go:

type IntegrationService struct {
	Service Service

	MockAPNS *apns.APNSMock
}

type buildTransactionFirestoreAdaptersDependencies struct {
	LoggerAdapter watermill.LoggerAdapter
}

var downloaderSet = wire.NewSet(app.NewDownloader)

var followChangePullerSet = wire.NewSet(app.NewFollowChangePuller)

var generatorSet = wire.NewSet(notifications.NewGenerator)
