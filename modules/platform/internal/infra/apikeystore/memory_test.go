package apikeystore

import (
	"context"
	"testing"
	"time"

	"github.com/Santozz-x/Aureon/modules/platform/internal/domain"
)

func TestMemory_SaveAndFindBySecretHash(t *testing.T) {
	store := NewMemory()
	key := domain.APIKey{ID: "key_1", AccountID: "acct_1", SecretHash: "hash-1", CreatedAt: time.Now().UTC()}

	if err := store.Save(context.Background(), key); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.FindBySecretHash(context.Background(), "hash-1")
	if err != nil {
		t.Fatalf("FindBySecretHash: %v", err)
	}
	if got.ID != key.ID {
		t.Fatalf("FindBySecretHash returned id %s, want %s", got.ID, key.ID)
	}
}

func TestMemory_FindBySecretHash_NotFound(t *testing.T) {
	store := NewMemory()

	_, err := store.FindBySecretHash(context.Background(), "does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown secret hash, got nil")
	}
}

func TestMemory_Revoke(t *testing.T) {
	store := NewMemory()
	key := domain.APIKey{ID: "key_1", AccountID: "acct_1", SecretHash: "hash-1", CreatedAt: time.Now().UTC()}
	if err := store.Save(context.Background(), key); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := store.Revoke(context.Background(), "key_1"); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	got, err := store.FindBySecretHash(context.Background(), "hash-1")
	if err != nil {
		t.Fatalf("FindBySecretHash: %v", err)
	}
	if !got.Revoked() {
		t.Fatal("expected key to be revoked")
	}
}

func TestMemory_Revoke_UnknownID(t *testing.T) {
	store := NewMemory()

	if err := store.Revoke(context.Background(), "does-not-exist"); err == nil {
		t.Fatal("expected error revoking unknown id, got nil")
	}
}
