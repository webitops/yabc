package blockchain

import "yabc/config"

type Config struct {
	config.Config
}

func NewBlockchainConfig(options map[string]string) *Config {
	return &Config{
		Config: config.Config{
			Options: options,
		},
	}
}

func (c *Config) getNodeAddress() string {
	value, exists := c.Get("node-address")
	if !exists {
		return ""
	}
	return value
}
