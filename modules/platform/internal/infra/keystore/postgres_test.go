package keystore

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	chainport "github.com/jeielsantos/aureon/modules/contracts"
	"github.com/jeielsantos/aureon/modules/platform/internal/infra/db"
)

// testDB connects to a real Postgres and applies migrations. Skipped
// unless TEST_DATABASE_URL is set — these are integration tests, not
// part of the default `make test` run.
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

func testEncryptionKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("generate encryption key: %v", err)
	}
	return key
}

func uniqueAddress(t *testing.T) chainport.Address {
	t.Helper()
	return chainport.Address(fmt.Sprintf("0xtest-%d", time.Now().UnixNano()))
}

func TestPostgres_PutAndGet(t *testing.T) {
	conn := testDB(t)
	store, err := NewPostgres(conn, testEncryptionKey(t))
	if err != nil {
		t.Fatalf("NewPostgres: %v", err)
	}

	address := uniqueAddress(t)
	privateKey := []byte("super-secret-private-key-bytes")

	if err := store.Put(context.Background(), address, privateKey); err != nil {
		t.Fatalf("Put: %v", err)
	}

	got, err := store.Get(context.Background(), address)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !bytes.Equal(got, privateKey) {
		t.Fatalf("Get = %q, want %q", got, privateKey)
	}
}

func TestPostgres_Get_NotFound(t *testing.T) {
	conn := testDB(t)
	store, err := NewPostgres(conn, testEncryptionKey(t))
	if err != nil {
		t.Fatalf("NewPostgres: %v", err)
	}

	_, err = store.Get(context.Background(), uniqueAddress(t))
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestPostgres_EncryptsAtRest(t *testing.T) {
	conn := testDB(t)
	store, err := NewPostgres(conn, testEncryptionKey(t))
	if err != nil {
		t.Fatalf("NewPostgres: %v", err)
	}

	address := uniqueAddress(t)
	privateKey := []byte("super-secret-private-key-bytes")
	if err := store.Put(context.Background(), address, privateKey); err != nil {
		t.Fatalf("Put: %v", err)
	}

	var raw []byte
	err = conn.QueryRow(`SELECT encrypted_private_key FROM wallet_keys WHERE address = $1`, string(address)).Scan(&raw)
	if err != nil {
		t.Fatalf("read raw column: %v", err)
	}
	if bytes.Contains(raw, privateKey) {
		t.Fatal("plaintext private key found in stored ciphertext")
	}
}

func TestPostgres_Get_WrongEncryptionKey(t *testing.T) {
	conn := testDB(t)
	store, err := NewPostgres(conn, testEncryptionKey(t))
	if err != nil {
		t.Fatalf("NewPostgres: %v", err)
	}

	address := uniqueAddress(t)
	if err := store.Put(context.Background(), address, []byte("secret")); err != nil {
		t.Fatalf("Put: %v", err)
	}

	otherStore, err := NewPostgres(conn, testEncryptionKey(t))
	if err != nil {
		t.Fatalf("NewPostgres: %v", err)
	}

	if _, err := otherStore.Get(context.Background(), address); err == nil {
		t.Fatal("expected decryption error with a different key, got nil")
	}
}
