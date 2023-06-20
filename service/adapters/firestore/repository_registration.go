package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type RegistrationRepository struct {
	tx *firestore.Transaction
}

func NewRegistrationRepository(tx *firestore.Transaction) *RegistrationRepository {
	return &RegistrationRepository{tx: tx}
}

func (r RegistrationRepository) Save(registration domain.Registration) error {
	return errors.New("save not implemented")
}
