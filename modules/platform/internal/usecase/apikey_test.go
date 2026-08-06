package usecase

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jeielsantos/aureon/modules/platform/internal/domain"
)

type fakeAPIKeyStore struct {
	mu   sync.Mutex
	keys map[domain.APIKeyID]domain.APIKey
}

func newFakeAPIKeyStore() *fakeAPIKeyStore {
	return &fakeAPIKeyStore{keys: make(map[domain.APIKeyID]domain.APIKey)}
}

func (f *fakeAPIKeyStore) Save(ctx context.Context, key domain.APIKey) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.keys[key.ID] = key
	return nil
}

func (f *fakeAPIKeyStore) FindBySecretHash(ctx context.Context, secretHash string) (domain.APIKey, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, key := range f.keys {
		if key.SecretHash == secretHash {
			return key, nil
		}
	}
	return domain.APIKey{}, fmt.Errorf("no matching key")
}

func (f *fakeAPIKeyStore) Revoke(ctx context.Context, id domain.APIKeyID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	key, ok := f.keys[id]
	if !ok {
		return fmt.Errorf("no key with id %q", id)
	}
	now := time.Now().UTC()
	key.RevokedAt = &now
	f.keys[id] = key
	return nil
}

func TestAPIKeyService_IssueAndAuthenticate(t *testing.T) {
	service := NewAPIKeyService(newFakeAPIKeyStore())

	id, secret, err := service.Issue(context.Background(), "acct_1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if id == "" || secret == "" {
		t.Fatal("Issue returned an empty id or secret")
	}

	key, err := service.Authenticate(context.Background(), secret)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if key.ID != id {
		t.Fatalf("Authenticate returned key %s, want %s", key.ID, id)
	}
	if key.AccountID != "acct_1" {
		t.Fatalf("Authenticate returned account %s, want acct_1", key.AccountID)
	}
}

func TestAPIKeyService_Authenticate_WrongSecret(t *testing.T) {
	service := NewAPIKeyService(newFakeAPIKeyStore())

	if _, _, err := service.Issue(context.Background(), "acct_1"); err != nil {
		t.Fatalf("Issue: %v", err)
	}

	_, err := service.Authenticate(context.Background(), "aur_not-the-real-secret")
	if err != ErrInvalidAPIKey {
		t.Fatalf("Authenticate error = %v, want ErrInvalidAPIKey", err)
	}
}

func TestAPIKeyService_Authenticate_AfterRevoke(t *testing.T) {
	service := NewAPIKeyService(newFakeAPIKeyStore())

	id, secret, err := service.Issue(context.Background(), "acct_1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	if err := service.Revoke(context.Background(), id); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	_, err = service.Authenticate(context.Background(), secret)
	if err != ErrInvalidAPIKey {
		t.Fatalf("Authenticate error = %v, want ErrInvalidAPIKey", err)
	}
}

func TestAPIKeyService_Issue_GeneratesDistinctSecrets(t *testing.T) {
	service := NewAPIKeyService(newFakeAPIKeyStore())

	_, first, err := service.Issue(context.Background(), "acct_1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	_, second, err := service.Issue(context.Background(), "acct_1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	if first == second {
		t.Fatalf("expected distinct secrets, got %s twice", first)
	}
}
