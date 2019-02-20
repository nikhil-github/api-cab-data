package wiring

import "time"

type Config struct {
	DBConfig DatabaseConfig
}

type DatabaseConfig struct {
	URL string

	Connections struct {
		Idle     int
		Lifetime time.Duration
		Max      int
	}
}
