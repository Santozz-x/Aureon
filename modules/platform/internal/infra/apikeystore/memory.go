// Package apikeystore provides usecase.APIKeyStore implementations.
package apikeystore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Santozz-x/Aureon/modules/platform/internal/domain"
)

// Memory is an in-process, non-persistent store. Keys are lost on
// restart. Interim until Sprint 5 lands Postgres — do not use in
// production. Mirrors the same caveat as infra/keystore.Memory (see
// TR-008 in docs/tradeoffs.md).
type Memory struct {
	mu   sync.RWMutex
	keys map[domain.APIKeyID]domain.APIKey
}

func NewMemory() *Memory {
	return &Memory{keys: make(map[domain.APIKeyID]domain.APIKey)}
}

func (m *Memory) Save(ctx context.Context, key domain.APIKey) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keys[key.ID] = key
	return nil
}

func (m *Memory) FindBySecretHash(ctx context.Context, secretHash string) (domain.APIKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, key := range m.keys {
		if key.SecretHash == secretHash {
			return key, nil
		}
	}
	return domain.APIKey{}, fmt.Errorf("apikeystore: no matching api key")
}

func (m *Memory) Revoke(ctx context.Context, id domain.APIKeyID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key, ok := m.keys[id]
	if !ok {
		return fmt.Errorf("apikeystore: no api key with id %q", id)
	}
	now := time.Now().UTC()
	key.RevokedAt = &now
	m.keys[id] = key
	return nil
}
