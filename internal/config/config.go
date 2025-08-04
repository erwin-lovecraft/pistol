package config

import (
	"os"
)

type Config struct {
	Port  string
	PGURL string
}

func ReadFromENV() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pgURL := os.Getenv("PG_URL")

	return Config{
		Port:  port,
		PGURL: pgURL,
	}
}
