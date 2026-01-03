package config

import "os"

func GetEnvFile() string {
	if os.Getenv("ENV") == "production" {
		return ".env.production"
	}
	return ".env"
}
