package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/watermill"
	watermillfirestore "github.com/ThreeDotsLabs/watermill-firestore/pkg/firestore"
)

const watermillRootCollection = "pubsub"

func NewWatermillPublisher(client *firestore.Client, logger watermill.LoggerAdapter) (*watermillfirestore.Publisher, error) {
	config := watermillfirestore.PublisherConfig{
		PubSubRootCollection:  watermillRootCollection,
		CustomFirestoreClient: client,
	}
	return watermillfirestore.NewPublisher(config, logger)
}

func NewWatermillSubscriber(client *firestore.Client, logger watermill.LoggerAdapter) (*watermillfirestore.Subscriber, error) {
	config := watermillfirestore.SubscriberConfig{
		PubSubRootCollection:  watermillRootCollection,
		CustomFirestoreClient: client,
	}
	return watermillfirestore.NewSubscriber(config, logger)
}
