package config

type NetworkConfig struct {
	Config
}

func NewNetworkConfig(options map[string]string) *NetworkConfig {
	return &NetworkConfig{
		Config: Config{
			options: options,
		},
	}
}

func (c *NetworkConfig) IsDebugEnabled() bool {
	value, exists := c.Get("debug")
	if !exists {
		return false
	}
	return value == "true"
}
