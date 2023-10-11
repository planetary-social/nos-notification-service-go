package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/config"
)

const (
	envPrefix = "NOTIFICATIONS"

	envNostrListenAddress              = "NOSTR_LISTEN_ADDRESS"
	envMetricsListenAddress            = "METRICS_LISTEN_ADDRESS"
	envFirestoreProjectID              = "FIRESTORE_PROJECT_ID"
	envFirestoreCredentialsJSONPath    = "FIRESTORE_CREDENTIALS_JSON_PATH"
	envAPNSTopic                       = "APNS_TOPIC"
	envAPNSCertificatePath             = "APNS_CERTIFICATE_PATH"
	envAPNSCertificatePassword         = "APNS_CERTIFICATE_PASSWORD"
	envEnvironment                     = "ENVIRONMENT"
	envLogLevel                        = "LOG_LEVEL"
	envGooglePubsubEnabled             = "GOOGLE_PUBSUB_ENABLED"
	envGooglePubsubProjectID           = "GOOGLE_PUBSUB_PROJECT_ID"
	envGooglePubsubCredentialsJSONPath = "GOOGLE_PUBSUB_CREDENTIALS_JSON_PATH"
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

	logLevel, err := c.loadLogLevel()
	if err != nil {
		return config.Config{}, errors.Wrap(err, "error loading the log level")
	}

	var firestoreCredentialsJSON []byte
	if p := c.getenv(envFirestoreCredentialsJSONPath); p != "" {
		f, err := os.Open(p)
		if err != nil {
			return config.Config{}, errors.Wrap(err, "error opening the firestore credentials file")
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return config.Config{}, errors.Wrap(err, "error reading the firestore credentials file")
		}

		firestoreCredentialsJSON = b
	}

	var googlePubSubCredentialsJSON []byte
	if p := c.getenv(envGooglePubsubCredentialsJSONPath); p != "" {
		f, err := os.Open(p)
		if err != nil {
			return config.Config{}, errors.Wrap(err, "error opening the pubsub credentials file")
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return config.Config{}, errors.Wrap(err, "error reading the pubsub credentials file")
		}

		googlePubSubCredentialsJSON = b
	}

	googlePubSubEnabled, err := c.getenvbool(envGooglePubsubEnabled)
	if err != nil {
		return config.Config{}, errors.Wrapf(err, "error loading variable '%s'", envGooglePubsubEnabled)
	}

	return config.NewConfig(
		c.getenv(envNostrListenAddress),
		c.getenv(envMetricsListenAddress),
		c.getenv(envFirestoreProjectID),
		firestoreCredentialsJSON,
		c.getenv(envAPNSTopic),
		c.getenv(envAPNSCertificatePath),
		c.getenv(envAPNSCertificatePassword),
		environment,
		logLevel,
		googlePubSubEnabled,
		c.getenv(envGooglePubsubProjectID),
		googlePubSubCredentialsJSON,
	)
}

func (c *EnvironmentConfigLoader) loadEnvironment() (config.Environment, error) {
	v := strings.ToUpper(c.getenv(envEnvironment))
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

func (c *EnvironmentConfigLoader) loadLogLevel() (logging.Level, error) {
	v := strings.ToUpper(c.getenv(envLogLevel))
	switch v {
	case "TRACE":
		return logging.LevelTrace, nil
	case "DEBUG":
		return logging.LevelDebug, nil
	case "ERROR":
		return logging.LevelError, nil
	case "DISABLED":
		return logging.LevelDisabled, nil
	case "":
		return logging.LevelDebug, nil
	default:
		return 0, fmt.Errorf("invalid log level requested '%s'", v)
	}
}

func (c *EnvironmentConfigLoader) getenv(key string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", envPrefix, key))
}

func (c *EnvironmentConfigLoader) getenvbool(key string) (bool, error) {
	switch v := strings.ToUpper(c.getenv(key)); v {
	case "":
		return false, nil
	case "TRUE":
		return true, nil
	case "FALSE":
		return false, nil
	default:
		return false, fmt.Errorf("unknow value '%s'", v)
	}
}
