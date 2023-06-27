package firestore

import (
	"context"
	"encoding/hex"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/iterator"
)

const (
	collectionTags             = "tags"
	collectionTagsValues       = "tags"
	collectionTagsValuesEvents = "events"
)

type TagRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewTagRepository(
	client *firestore.Client,
	tx *firestore.Transaction,
) *TagRepository {
	return &TagRepository{
		client: client,
		tx:     tx,
	}
}

func (e *TagRepository) Save(event domain.Event) error {
	for _, tag := range event.Tags() {
		keyTag := encodeStringAsHex(tag.Name().String())
		keyValue := encodeStringAsHex(tag.FirstValue())

		tagDocPath := e.client.Collection(collectionTags).Doc(keyTag)
		if err := e.tx.Set(tagDocPath, map[string]any{}, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error updating the tag doc")
		}

		tagValueDocPath := e.client.Collection(collectionTags).Doc(keyTag).Collection(collectionTagsValues).Doc(keyValue)
		if err := e.tx.Set(tagValueDocPath, map[string]any{}, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error updating the value doc")
		}

		tagValueEventDocPath := e.client.Collection(collectionTags).Doc(keyTag).Collection(collectionTagsValues).Doc(keyValue).Collection(collectionTagsValuesEvents).Doc(event.Id().Hex())
		tagValueEventDocData := map[string]any{
			eventFieldId:        event.Id().Hex(),
			eventFieldPublicKey: event.PubKey().Hex(),
			eventFieldCreatedAt: event.CreatedAt(),
			eventFieldKind:      event.Kind().Int(),
			eventFieldRaw:       event.Raw(),
		}
		if err := e.tx.Set(tagValueEventDocPath, tagValueEventDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error updating the event doc")
		}
	}
	return nil
}

func (e *TagRepository) GetEvents(ctx context.Context, name domain.EventTagName, value string, since, until *time.Time, limit int, events map[string]domain.Event) error {
	keyTag := encodeStringAsHex(name.String())
	keyValue := encodeStringAsHex(value)

	query := e.client.Collection(collectionTags).Doc(keyTag).Collection(collectionTagsValues).Doc(keyValue).Collection(collectionTagsValuesEvents).Query

	if since != nil {
		query = query.Where(eventFieldCreatedAt, ">", since)
	}

	if until != nil {
		query = query.Where(eventFieldCreatedAt, "<", until)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	docs := query.Documents(ctx)

	for {
		doc, err := docs.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return errors.Wrap(err, "error getting next document")
		}

		data := make(map[string]any)
		if err := doc.DataTo(&data); err != nil {
			return errors.Wrap(err, "error reading document data")
		}

		event, err := domain.NewEventFromRaw(data[eventFieldRaw].([]byte))
		if err != nil {
			return errors.Wrap(err, "error creating the event")
		}

		events[event.Id().Hex()] = event
	}
	return nil
}

func encodeStringAsHex(s string) string {
	return hex.EncodeToString([]byte(s))
}
