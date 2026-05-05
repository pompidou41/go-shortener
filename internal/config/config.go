package config

import "os"

type Config struct {
	SecretSalt string
}

func New() *Config {
	return &Config{
		SecretSalt: string(getEnv("SECRET_SALT", "")),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
