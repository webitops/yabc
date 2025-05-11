package config

type Config struct {
	options map[string]string
}

func (c *Config) Get(key string) (string, bool) {
	element, exists := c.options[key]
	return element, exists
}

func (c *Config) Set(key string, value string) {
	c.options[key] = value
}
