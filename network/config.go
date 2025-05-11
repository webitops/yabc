package network

import "yabc/config"

type Config struct {
	config.Config
}

func NewNetworkConfig(options map[string]string) *Config {
	return &Config{
		Config: config.Config{
			Options: options,
		},
	}
}

func (c *Config) IsDebugEnabled() bool {
	value, exists := c.Get("debug")
	if !exists {
		return false
	}
	return value == "true"
}
