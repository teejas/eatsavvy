package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

type DatabaseClient struct {
	Db  *pgx.Conn
	Ctx context.Context
}

func NewDatabaseClient() *DatabaseClient {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, generateConnectionString())
	if err != nil {
		slog.Error("[db.NewDatabaseClient] Failed to connect to database", "error", err)
		return nil
	}
	slog.Info("[db.NewDatabaseClient] Connected to database")
	return &DatabaseClient{Db: db, Ctx: ctx}
}

func (dc *DatabaseClient) Close() {
	if dc.Db == nil {
		slog.Error("[db.Close] Database connection is nil")
		return
	}
	err := dc.Db.Close(dc.Ctx)
	if err != nil {
		slog.Error("[db.Close] Failed to close database connection", "error", err)
	}
	slog.Info("[db.Close] Closed database connection")
}

func generateConnectionString() string {
	// postgresql://[DB_USER]:[DB_PASSWORD]@[DB_HOST]:[DB_PORT]/[DB_NAME]
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	slog.Info("[db.generateConnectionString] Generated connection string", "connString", connString)
	return connString
}
