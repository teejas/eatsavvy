package main

import (
	"eatsavvy/internal/config"
	"eatsavvy/internal/worker"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(config.GetEnvFile())
	if err != nil {
		slog.Error("[worker.main] Failed to load .env file", "error", err)
	}
	worker := worker.NewWorker()
	worker.Start()
}
