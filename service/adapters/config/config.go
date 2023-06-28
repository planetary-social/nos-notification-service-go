package config

import (
	"fmt"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/config"
	"os"
	"strings"
)

const (
	envPrefix = "NOTIFICATIONS"

	envNostrListenAddress      = "NOSTR_LISTEN_ADDRESS"
	envFirestoreProjectID      = "FIRESTORE_PROJECT_ID"
	envAPNSTopic               = "APNS_TOPIC"
	envAPNSCertificatePath     = "APNS_CERTIFICATE_PATH"
	envAPNSCertificatePassword = "APNS_CERTIFICATE_PASSWORD"
	envEnvironment             = "ENVIRONMENT"
)

type EnvironmentConfigLoader struct {
}

func NewEnvironmentConfigLoader() *EnvironmentConfigLoader {
	return &EnvironmentConfigLoader{}
}

func (c *EnvironmentConfigLoader) Load() (config.Config, error) {
	environment, err := c.loadEnvironment()
	if err != nil {
		return config.Config{}, errors.Wrap(err, "error loading the environment setting")
	}

	return config.NewConfig(
		c.getenv(envNostrListenAddress),
		c.getenv(envFirestoreProjectID),
		c.getenv(envAPNSTopic),
		c.getenv(envAPNSCertificatePath),
		c.getenv(envAPNSCertificatePassword),
		environment,
	)
}

func (c *EnvironmentConfigLoader) loadEnvironment() (config.Environment, error) {
	v := strings.ToUpper(os.Getenv(envEnvironment))
	switch v {
	case "PRODUCTION":
		return config.EnvironmentProduction, nil
	case "DEVELOPMENT":
		return config.EnvironmentDevelopment, nil
	case "":
		return config.EnvironmentProduction, nil
	default:
		return config.Environment{}, fmt.Errorf("invalid environment requested '%s'", v)
	}
}

func (c *EnvironmentConfigLoader) getenv(key string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", envPrefix, key))
}
