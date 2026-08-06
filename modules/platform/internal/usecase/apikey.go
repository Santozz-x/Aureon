package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jeielsantos/aureon/modules/platform/internal/domain"
)

var ErrInvalidAPIKey = errors.New("usecase: invalid or revoked api key")

type APIKeyStore interface {
	Save(ctx context.Context, key domain.APIKey) error
	FindBySecretHash(ctx context.Context, secretHash string) (domain.APIKey, error)
	Revoke(ctx context.Context, id domain.APIKeyID) error
}

type APIKeyService struct {
	store APIKeyStore
}

func NewAPIKeyService(store APIKeyStore) *APIKeyService {
	return &APIKeyService{store: store}
}

// Issue returns the key's id (safe to log/display) and its raw secret
// (shown to the caller exactly once — only its hash is ever stored).
func (s *APIKeyService) Issue(ctx context.Context, accountID domain.AccountID) (domain.APIKeyID, string, error) {
	secret, err := randomToken("aur_")
	if err != nil {
		return "", "", fmt.Errorf("usecase: generate api key secret: %w", err)
	}
	rawID, err := randomToken("")
	if err != nil {
		return "", "", fmt.Errorf("usecase: generate api key id: %w", err)
	}
	id := domain.APIKeyID(rawID)

	key := domain.APIKey{
		ID:         id,
		AccountID:  accountID,
		SecretHash: hashSecret(secret),
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.store.Save(ctx, key); err != nil {
		return "", "", fmt.Errorf("usecase: save api key: %w", err)
	}

	return id, secret, nil
}

func (s *APIKeyService) Revoke(ctx context.Context, id domain.APIKeyID) error {
	if err := s.store.Revoke(ctx, id); err != nil {
		return fmt.Errorf("usecase: revoke api key %s: %w", id, err)
	}
	return nil
}

func (s *APIKeyService) Authenticate(ctx context.Context, secret string) (domain.APIKey, error) {
	key, err := s.store.FindBySecretHash(ctx, hashSecret(secret))
	if err != nil || key.Revoked() {
		return domain.APIKey{}, ErrInvalidAPIKey
	}
	return key, nil
}

func randomToken(prefix string) (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(buf), nil
}

// hashSecret uses SHA-256, not a slow password hash (bcrypt/argon2): API
// keys are high-entropy random tokens, not user-chosen passwords, so
// there's no low-entropy search space for a slow hash to protect against.
func hashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
