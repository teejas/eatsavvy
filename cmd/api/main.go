package main

import (
	"eatsavvy/internal/api"
	"eatsavvy/pkg/utils"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(utils.GetEnvFile())
	if err != nil {
		slog.Error("Failed to load .env file", "error", err)
	}

	api.StartServer("8080")
}
