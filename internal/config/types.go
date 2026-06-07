package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Route struct {
	Method   string   `yaml:"method"`
	Path     string   `yaml:"path"`
	Response Response `yaml:"response"`
}

type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

func (c *Config) setDefaults() {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	for i := range c.Routes {
		if c.Routes[i].Response.Status == 0 {
			c.Routes[i].Response.Status = 200
		}
	}
}
