package domain

import (
	"encoding/json"
	"net"
	"net/url"
	"strings"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal"
)

// todo make sure that the registration was sent by one of those public keys?
type Registration struct {
	apnsToken APNSToken
	publicKey PublicKey
	relays    []RelayAddress
}

func NewRegistrationFromEvent(event Event) (Registration, error) {
	var v registrationTransport
	if err := json.Unmarshal([]byte(event.Content()), &v); err != nil {
		return Registration{}, errors.Wrap(err, "error unmarshaling content")
	}

	apnsToken, err := NewAPNSTokenFromHex(v.APNSToken)
	if err != nil {
		return Registration{}, errors.Wrap(err, "error creating an apns token")
	}

	publicKey, err := NewPublicKeyFromHex(v.PublicKey)
	if err != nil {
		return Registration{}, errors.Wrap(err, "error creating a public key")
	}

	relays, err := newRelays(v)
	if err != nil {
		return Registration{}, errors.Wrap(err, "error creating relay addresses")
	}

	if event.PubKey() != publicKey {
		return Registration{}, errors.New("public key doesn't match public key from event")
	}

	return Registration{
		apnsToken: apnsToken,
		publicKey: publicKey,
		relays:    relays,
	}, nil
}

func newRelays(v registrationTransport) ([]RelayAddress, error) {
	var relays []RelayAddress
	for _, relayTransport := range v.Relays {
		address, err := NewRelayAddress(relayTransport.Address)
		if err != nil {
			return nil, errors.Wrap(err, "error creating relay address")
		}
		relays = append(relays, address)
	}

	if len(relays) == 0 {
		return nil, errors.New("missing relays")
	}

	return relays, nil
}

func (r Registration) APNSToken() APNSToken {
	return r.apnsToken
}

func (p Registration) PublicKey() PublicKey {
	return p.publicKey
}

func (p Registration) Relays() []RelayAddress {
	return internal.CopySlice(p.relays)
}

type RelayAddress struct {
	s               string
	hostWithoutPort string
}

func NewRelayAddress(s string) (RelayAddress, error) {
	if !strings.HasPrefix(s, "ws://") && !strings.HasPrefix(s, "wss://") {
		return RelayAddress{}, errors.New("invalid protocol")
	}

	s = strings.TrimSpace(s)
	s = strings.TrimRight(s, "/")

	u, err := url.Parse(s)
	if err != nil {
		return RelayAddress{}, errors.Wrap(err, "url parse error")
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		return RelayAddress{}, errors.New("invalid protocol")
	}

	u.Host = strings.ToLower(u.Host)
	hostWithoutPort, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		hostWithoutPort = u.Host
	}
	normalizedURI := u.String()

	return RelayAddress{
		s:               normalizedURI,
		hostWithoutPort: hostWithoutPort,
	}, nil
}

func (r RelayAddress) String() string {
	return r.s
}

func (r RelayAddress) HostWithoutPort() string {
	return r.hostWithoutPort
}

type registrationTransport struct {
	APNSToken string           `json:"apnsToken"`
	PublicKey string           `json:"publicKey"`
	Relays    []relayTransport `json:"relays"`
}

type relayTransport struct {
	Address string `json:"address"`
}
