package keystore

import (
	"context"
	"fmt"
	"sync"

	chainport "github.com/Santozz-x/Aureon/modules/contracts"
)

// Memory is an in-process, non-persistent chainport.KeyStore. Keys are lost
// on restart and never touch disk. It exists only to unblock Wallet/Transaction
// API development before Sprint 5 lands a real (encrypted, persistent) store —
// it must never hold keys for funds with real value. See TR-008 in
// docs/tradeoffs.md.
type Memory struct {
	mu   sync.RWMutex
	keys map[chainport.Address][]byte
}

func NewMemory() *Memory {
	return &Memory{keys: make(map[chainport.Address][]byte)}
}

func (m *Memory) Put(ctx context.Context, address chainport.Address, privateKey []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keys[address] = privateKey
	return nil
}

func (m *Memory) Get(ctx context.Context, address chainport.Address) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key, ok := m.keys[address]
	if !ok {
		return nil, fmt.Errorf("keystore: no key stored for address %q", address)
	}
	return key, nil
}
