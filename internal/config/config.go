package config

import (
	"os"
)

type Config struct {
	Port string
}

func ReadFromENV() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
	}
}
