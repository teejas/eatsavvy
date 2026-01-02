package main

import (
	"eatsavvy/internal/worker"
	"eatsavvy/pkg/utils"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(utils.GetEnvFile())
	if err != nil {
		slog.Error("Failed to load .env file", "error", err)
	}
	worker := worker.NewWorker()
	worker.Start()
}
