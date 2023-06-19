package config

type Config struct {
	nostrListenAddress string
}

// NewConfig creates a new config with the following options:
//
// nostrListenAddress:
//
//	listen address for the websocket connections in the format accepted by the
//	standard library.
//
//	Optional, defaults to ":8008".
func NewConfig(nostrListenAddress string) Config {
	c := Config{
		nostrListenAddress: nostrListenAddress,
	}

	c.setDefaults()
	return c
}

func (c *Config) NostrListenAddress() string {
	return c.nostrListenAddress
}

func (c *Config) setDefaults() {
	if c.nostrListenAddress == "" {
		c.nostrListenAddress = ":8008"
	}
}
