package config

import "os"

type Config struct {
	Port    string
	AppEnv  string
	GinMode string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8313"
	}
	return &Config{
		Port:    port,
		AppEnv:  getEnv("APP_ENV", "development"),
		GinMode: getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
