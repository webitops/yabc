package config

type Config struct {
	Options map[string]string
}

func (c *Config) Get(key string) (string, bool) {
	element, exists := c.Options[key]
	return element, exists
}

func (c *Config) Set(key string, value string) {
	c.Options[key] = value
}
