package config

import "os"

type Config struct {
	Storage string
	Port    string
}

func New() *Config {
	return &Config{
		Storage: string(getEnv("STORAGE", "memory")),
		Port:    string(getEnv("PORT", "3000")),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
