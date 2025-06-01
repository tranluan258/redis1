package config

import "fmt"

type Config struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Password string `json:"password,omitempty"`
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if len(c.Password) > 0 && len(c.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}
