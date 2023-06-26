package config

import (
	"fmt"

	"github.com/boreq/errors"
)

type Environment struct {
	s string
}

var (
	EnvironmentProduction  = Environment{"production"}
	EnvironmentDevelopment = Environment{"development"}
)

type Config struct {
	nostrListenAddress string
	firestoreProjectID string

	apnsTopic               string
	apnsCertificatePath     string
	apnsCertificatePassword string

	environment Environment
}

// NewConfig creates a new config with the following options:
//
// - nostrListenAddress:
//
//	Listen address for the websocket connections in the format accepted by the
//	standard library.
//
//	Optional, defaults to ":8008" if empty.
//
// - firestoreProjectID
//
//	Your firestore project id.
//
//	Required.
//
// - apnsTopic
//
//	Topic on which APNs notifications will be sent.
//
//	Required.
//
// - apnsCertificatePath
//
//	Path to your APNs certificate.
//
//	Required.
func NewConfig(
	nostrListenAddress string,
	firestoreProjectID string,
	apnsTopic string,
	apnsCertificatePath string,
	apnsCertificatePassword string,
	environment Environment,
) (Config, error) {
	c := Config{
		nostrListenAddress:      nostrListenAddress,
		firestoreProjectID:      firestoreProjectID,
		apnsTopic:               apnsTopic,
		apnsCertificatePath:     apnsCertificatePath,
		apnsCertificatePassword: apnsCertificatePassword,
		environment:             environment,
	}

	c.setDefaults()
	if err := c.validate(); err != nil {
		return Config{}, errors.Wrap(err, "invalid config")
	}

	return c, nil
}

func (c *Config) NostrListenAddress() string {
	return c.nostrListenAddress
}

func (c *Config) FirestoreProjectID() string {
	return c.firestoreProjectID
}

func (c *Config) APNSTopic() string {
	return c.apnsTopic
}

func (c *Config) APNSCertificatePath() string {
	return c.apnsCertificatePath
}

func (c *Config) APNSCertificatePassword() string {
	return c.apnsCertificatePassword
}

func (c *Config) Environment() Environment {
	return c.environment
}

func (c *Config) setDefaults() {
	if c.nostrListenAddress == "" {
		c.nostrListenAddress = ":8008"
	}
}

func (c *Config) validate() error {
	if c.nostrListenAddress == "" {
		return errors.New("missing nostr listen address")
	}

	if c.firestoreProjectID == "" {
		return errors.New("missing firestore project id")
	}

	if c.apnsTopic == "" {
		return errors.New("missing APNs topic")
	}

	if c.apnsCertificatePath == "" {
		return errors.New("missing APNs certificate path")
	}

	switch c.environment {
	case EnvironmentProduction:
	case EnvironmentDevelopment:
	default:
		return fmt.Errorf("unknown environment '%+v'", c.environment)
	}

	return nil
}
