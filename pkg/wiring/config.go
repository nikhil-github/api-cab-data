package wiring

import "time"

// Config wraps app configs.
type Config struct {
	DB   DBConfig
	HTTP struct {
		Port int `envconfig:"default=3000"`
	}
	LOG struct {
		Level string
	}
}

// DBConfig wraps DB configs.
type DBConfig struct {
	URL         string
	Connections struct {
		Idle     int           `envconfig:"default=10"`
		Lifetime time.Duration `envconfig:"default=5m"`
		Max      int           `envconfig:"default=20"`
	}
}
