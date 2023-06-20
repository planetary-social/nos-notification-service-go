package domain

import (
	"encoding/json"

	"github.com/boreq/errors"
)

type Registration struct {
	publicKeys []PublicKeyWithRelays
	locale     Locale
	apnsToken  APNSToken
}

func NewRegistrationFromEvent(event Event) (Registration, error) {
	var v registrationTransport
	if err := json.Unmarshal(event.Content(), &v); err != nil {
		return Registration{}, errors.Wrap(err, "error unmarshaling content")
	}

	publicKeysWithRelays, err := newPublicKeysWithRelays(v)
	if err != nil {
		return Registration{}, errors.Wrap(err, "error creating public keys with relays")
	}

	locale, err := NewLocale(v.Locale)
	if err != nil {
		return Registration{}, errors.Wrap(err, "error creating locale")
	}

	apnsToken, err := NewAPNSToken(v.APNSToken)
	if err != nil {
		return Registration{}, errors.Wrap(err, "error creating APNS token")
	}

	return Registration{
		publicKeys: publicKeysWithRelays,
		locale:     locale,
		apnsToken:  apnsToken,
	}, nil
}

func newPublicKeysWithRelays(v registrationTransport) ([]PublicKeyWithRelays, error) {
	var publicKeysWithRelays []PublicKeyWithRelays
	for _, publicKeyWithRelaysTransport := range v.PublicKeys {
		publicKey, err := NewPublicKey(publicKeyWithRelaysTransport.PublicKey)
		if err != nil {
			return nil, errors.Wrap(err, "error creating public key")
		}

		var relays []RelayAddress
		for _, relayTransport := range publicKeyWithRelaysTransport.Relays {
			address, err := NewRelayAddress(relayTransport.Address)
			if err != nil {
				return nil, errors.Wrap(err, "error creating relay address")
			}
			relays = append(relays, address)
		}

		v, err := NewPublicKeyWithRelays(publicKey, relays)
		if err != nil {
			return nil, errors.Wrap(err, "error creating public key with relays")
		}

		publicKeysWithRelays = append(publicKeysWithRelays, v)
	}

	if len(publicKeysWithRelays) == 0 {
		return nil, errors.New("empty public keys with relay")
	}

	return publicKeysWithRelays, nil
}

type PublicKeyWithRelays struct {
	publicKey PublicKey
	relays    []RelayAddress
}

func NewPublicKeyWithRelays(publicKey PublicKey, relays []RelayAddress) (PublicKeyWithRelays, error) {
	// todo validate e.g. relays can't be empty

	return PublicKeyWithRelays{publicKey: publicKey, relays: relays}, nil
}

type RelayAddress struct {
	s string
}

func NewRelayAddress(s string) (RelayAddress, error) {
	// todo validate

	return RelayAddress{s: s}, nil
}

type registrationTransport struct {
	PublicKeys []publicKeysWithRelaysTransport `json:"publicKeys"`
	Locale     string                          `json:"locale"`
	APNSToken  string                          `json:"apnsToken"`
}

type publicKeysWithRelaysTransport struct {
	PublicKey string           `json:"publicKey"`
	Relays    []relayTransport `json:"relays"`
}

type relayTransport struct {
	Address string `json:"address"`
}
