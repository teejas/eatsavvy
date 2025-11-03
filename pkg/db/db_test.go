package db

import (
	"os"
	"testing"
)

func TestGenerateConnectionString(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_PORT", "5432")
	connString := generateConnectionString()
	if connString != "postgresql://postgres:password@localhost:5432/postgres" {
		t.Errorf("Expected connection string to be 'postgresql://postgres:password@localhost:5432/postgres', but got '%s'", connString)
	}
}
