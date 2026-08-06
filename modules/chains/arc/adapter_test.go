package arc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/jeielsantos/aureon/modules/chains/arc/rpc"
	chainport "github.com/jeielsantos/aureon/modules/contracts"
)

type fakeKeyStore struct {
	mu   sync.Mutex
	keys map[chainport.Address][]byte
}

func newFakeKeyStore() *fakeKeyStore {
	return &fakeKeyStore{keys: make(map[chainport.Address][]byte)}
}

func (f *fakeKeyStore) Put(ctx context.Context, address chainport.Address, privateKey []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.keys[address] = privateKey
	return nil
}

func (f *fakeKeyStore) Get(ctx context.Context, address chainport.Address) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key, ok := f.keys[address]
	if !ok {
		return nil, fmt.Errorf("no key for %s", address)
	}
	return key, nil
}

func TestAdapter_CreateWallet(t *testing.T) {
	keys := newFakeKeyStore()
	adapter := NewAdapter(nil, keys)

	address, err := adapter.CreateWallet(context.Background())
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}
	if address == "" {
		t.Fatal("CreateWallet returned an empty address")
	}

	storedKey, err := keys.Get(context.Background(), address)
	if err != nil {
		t.Fatalf("expected a stored key for %s: %v", address, err)
	}

	privateKey, err := crypto.ToECDSA(storedKey)
	if err != nil {
		t.Fatalf("stored key is not a valid ECDSA private key: %v", err)
	}
	gotAddress := chainport.Address(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	if gotAddress != address {
		t.Fatalf("stored key derives to address %s, want %s", gotAddress, address)
	}
}

func TestAdapter_CreateWallet_GeneratesDistinctWallets(t *testing.T) {
	keys := newFakeKeyStore()
	adapter := NewAdapter(nil, keys)

	first, err := adapter.CreateWallet(context.Background())
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}
	second, err := adapter.CreateWallet(context.Background())
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	if first == second {
		t.Fatalf("expected distinct addresses, got %s twice", first)
	}
}

type rpcRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
}

func newMockRPCServer(t *testing.T, result interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: result}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
}

func TestAdapter_GetBalance(t *testing.T) {
	want := big.NewInt(42_000_000)
	srv := newMockRPCServer(t, fmt.Sprintf("0x%x", want))
	defer srv.Close()

	client, err := rpc.Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	adapter := NewAdapter(client, newFakeKeyStore())

	balance, err := adapter.GetBalance(context.Background(), "0x1234567890123456789012345678901234567890")
	if err != nil {
		t.Fatalf("GetBalance: %v", err)
	}
	if balance.Amount != want.String() {
		t.Fatalf("GetBalance amount = %s, want %s", balance.Amount, want.String())
	}
	if balance.Symbol != "USDC" {
		t.Fatalf("GetBalance symbol = %s, want USDC", balance.Symbol)
	}
}

func TestAdapter_GetBalance_InvalidAddress(t *testing.T) {
	adapter := NewAdapter(nil, newFakeKeyStore())

	_, err := adapter.GetBalance(context.Background(), "not-an-address")
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
}
