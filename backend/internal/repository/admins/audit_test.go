package admins

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestGetAuditLogsScan(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5439/telegramshop?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("database not available: %v", err)
	}

	repo := NewRepository(db)
	logs, err := repo.GetAuditLogs(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetAuditLogs failed: %v", err)
	}
	t.Logf("got %d logs", len(logs))
}
