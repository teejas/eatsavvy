package main

import (
	"eatsavvy/internal/api"
	"eatsavvy/internal/config"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(config.GetEnvFile())
	if err != nil {
		slog.Error("[api.main] Failed to load .env file", "error", err)
	}

	api.StartServer("8080")
}
