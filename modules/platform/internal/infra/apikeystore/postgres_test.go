package apikeystore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Santozz-x/Aureon/modules/platform/internal/domain"
	"github.com/Santozz-x/Aureon/modules/platform/internal/infra/db"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping Postgres integration test")
	}

	conn, err := db.Open(context.Background(), url)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	return conn
}

func uniqueID(t *testing.T) domain.APIKeyID {
	t.Helper()
	return domain.APIKeyID(fmt.Sprintf("key_test_%d", time.Now().UnixNano()))
}

func TestPostgres_SaveAndFindBySecretHash(t *testing.T) {
	conn := testDB(t)
	store := NewPostgres(conn)

	key := domain.APIKey{
		ID:         uniqueID(t),
		AccountID:  "acct_1",
		SecretHash: fmt.Sprintf("hash-%d", time.Now().UnixNano()),
		CreatedAt:  time.Now().UTC().Truncate(time.Microsecond),
	}

	if err := store.Save(context.Background(), key); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.FindBySecretHash(context.Background(), key.SecretHash)
	if err != nil {
		t.Fatalf("FindBySecretHash: %v", err)
	}
	if got.ID != key.ID || got.AccountID != key.AccountID {
		t.Fatalf("FindBySecretHash = %+v, want id=%s account=%s", got, key.ID, key.AccountID)
	}
	if got.Revoked() {
		t.Fatal("expected freshly saved key to not be revoked")
	}
}

func TestPostgres_FindBySecretHash_NotFound(t *testing.T) {
	conn := testDB(t)
	store := NewPostgres(conn)

	_, err := store.FindBySecretHash(context.Background(), "does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown secret hash, got nil")
	}
}

func TestPostgres_Revoke(t *testing.T) {
	conn := testDB(t)
	store := NewPostgres(conn)

	key := domain.APIKey{
		ID:         uniqueID(t),
		AccountID:  "acct_1",
		SecretHash: fmt.Sprintf("hash-%d", time.Now().UnixNano()),
		CreatedAt:  time.Now().UTC().Truncate(time.Microsecond),
	}
	if err := store.Save(context.Background(), key); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := store.Revoke(context.Background(), key.ID); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	got, err := store.FindBySecretHash(context.Background(), key.SecretHash)
	if err != nil {
		t.Fatalf("FindBySecretHash: %v", err)
	}
	if !got.Revoked() {
		t.Fatal("expected key to be revoked")
	}
}

func TestPostgres_Revoke_UnknownID(t *testing.T) {
	conn := testDB(t)
	store := NewPostgres(conn)

	if err := store.Revoke(context.Background(), uniqueID(t)); err == nil {
		t.Fatal("expected error revoking unknown id, got nil")
	}
}
