// End-to-end tests of the Gateway's HTTP surface: real router, real
// usecase services, real middleware — only the blockchain network itself
// is faked, since no integration test should depend on a live testnet.
package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	chainport "github.com/Santozz-x/Aureon/modules/contracts"
	"github.com/Santozz-x/Aureon/modules/platform/internal/adapter/rest"
	"github.com/Santozz-x/Aureon/modules/platform/internal/adapter/rest/middleware"
	"github.com/Santozz-x/Aureon/modules/platform/internal/infra/apikeystore"
	"github.com/Santozz-x/Aureon/modules/platform/internal/usecase"
)

type fakeAdapter struct {
	network chainport.Network
}

func (f *fakeAdapter) Network() chainport.Network { return f.network }

func (f *fakeAdapter) CreateWallet(ctx context.Context) (chainport.Address, error) {
	return chainport.Address(fmt.Sprintf("0xfake-%d", time.Now().UnixNano())), nil
}

func (f *fakeAdapter) GetBalance(ctx context.Context, address chainport.Address) (chainport.Balance, error) {
	if address == "0xunknown" {
		return chainport.Balance{}, fmt.Errorf("fake: unknown address")
	}
	return chainport.Balance{Amount: "1000", Symbol: "TEST"}, nil
}

func (f *fakeAdapter) SendTransaction(ctx context.Context, tx chainport.Transaction) (chainport.TxHash, error) {
	return chainport.TxHash("0xfaketxhash"), nil
}

func (f *fakeAdapter) EstimateGas(ctx context.Context, tx chainport.Transaction) (uint64, error) {
	return 21000, nil
}

// newTestGateway wires the same components as cmd/gateway/main.go — real
// router, real usecase services, real auth/logging middleware — against a
// fakeAdapter instead of a live chain RPC connection.
func newTestGateway(t *testing.T) (*httptest.Server, *usecase.APIKeyService) {
	t.Helper()

	adapters := map[chainport.Network]chainport.Adapter{
		"testnet": &fakeAdapter{network: "testnet"},
	}

	walletService := usecase.NewWalletService(adapters)
	walletHandler := rest.NewWalletHandler(walletService)

	transactionService := usecase.NewTransactionService(adapters)
	transactionHandler := rest.NewTransactionHandler(transactionService)

	apiKeyService := usecase.NewAPIKeyService(apikeystore.NewMemory())
	apiKeyHandler := rest.NewAPIKeyHandler(apiKeyService)
	protect := middleware.RequireAPIKey(apiKeyService)

	router := rest.NewRouter(walletHandler, transactionHandler, apiKeyHandler, protect)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := middleware.Logging(logger)(router)

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return srv, apiKeyService
}

func issueTestAPIKey(t *testing.T, srv *httptest.Server) string {
	t.Helper()

	resp, err := srv.Client().Post(srv.URL+"/v1/apikeys", "application/json", strings.NewReader(`{"account_id":"acct_test"}`))
	if err != nil {
		t.Fatalf("issue api key: %v", err)
	}
	defer resp.Body.Close()

	var body struct {
		Secret string `json:"secret"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode issue response: %v", err)
	}
	return body.Secret
}

func TestGateway_Health_NoAuthRequired(t *testing.T) {
	srv, _ := newTestGateway(t)

	resp, err := srv.Client().Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestGateway_WalletRoutes_RequireAuth(t *testing.T) {
	srv, _ := newTestGateway(t)

	resp, err := srv.Client().Get(srv.URL + "/v1/testnet/wallets/0xabc/balance")
	if err != nil {
		t.Fatalf("GET balance: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestGateway_CreateWalletAndGetBalance(t *testing.T) {
	srv, _ := newTestGateway(t)
	secret := issueTestAPIKey(t, srv)

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/testnet/wallets", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("create wallet: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create wallet status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var created struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create wallet response: %v", err)
	}
	if created.Address == "" {
		t.Fatal("expected a non-empty address")
	}

	balReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/v1/testnet/wallets/"+created.Address+"/balance", nil)
	balReq.Header.Set("Authorization", "Bearer "+secret)
	balResp, err := srv.Client().Do(balReq)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}
	defer balResp.Body.Close()

	if balResp.StatusCode != http.StatusOK {
		t.Fatalf("balance status = %d, want %d", balResp.StatusCode, http.StatusOK)
	}

	var balance struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	}
	if err := json.NewDecoder(balResp.Body).Decode(&balance); err != nil {
		t.Fatalf("decode balance response: %v", err)
	}
	if balance.Amount != "1000" || balance.Symbol != "TEST" {
		t.Fatalf("balance = %+v, want amount=1000 symbol=TEST", balance)
	}
}

func TestGateway_GetBalance_UnknownAddress(t *testing.T) {
	srv, _ := newTestGateway(t)
	secret := issueTestAPIKey(t, srv)

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/v1/testnet/wallets/0xunknown/balance", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestGateway_UnsupportedNetwork(t *testing.T) {
	srv, _ := newTestGateway(t)
	secret := issueTestAPIKey(t, srv)

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/does-not-exist/wallets", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("create wallet: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestGateway_SendTransactionAndEstimateGas(t *testing.T) {
	srv, _ := newTestGateway(t)
	secret := issueTestAPIKey(t, srv)

	body := `{"from":"0xfrom","to":"0xto","value":"1000"}`

	sendReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/testnet/transactions", strings.NewReader(body))
	sendReq.Header.Set("Authorization", "Bearer "+secret)
	sendReq.Header.Set("Content-Type", "application/json")
	sendResp, err := srv.Client().Do(sendReq)
	if err != nil {
		t.Fatalf("send transaction: %v", err)
	}
	defer sendResp.Body.Close()

	if sendResp.StatusCode != http.StatusCreated {
		t.Fatalf("send transaction status = %d, want %d", sendResp.StatusCode, http.StatusCreated)
	}

	var sent struct {
		TxHash string `json:"tx_hash"`
	}
	if err := json.NewDecoder(sendResp.Body).Decode(&sent); err != nil {
		t.Fatalf("decode send response: %v", err)
	}
	if sent.TxHash != "0xfaketxhash" {
		t.Fatalf("tx_hash = %s, want 0xfaketxhash", sent.TxHash)
	}

	estReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/testnet/transactions/estimate", strings.NewReader(body))
	estReq.Header.Set("Authorization", "Bearer "+secret)
	estReq.Header.Set("Content-Type", "application/json")
	estResp, err := srv.Client().Do(estReq)
	if err != nil {
		t.Fatalf("estimate gas: %v", err)
	}
	defer estResp.Body.Close()

	if estResp.StatusCode != http.StatusOK {
		t.Fatalf("estimate gas status = %d, want %d", estResp.StatusCode, http.StatusOK)
	}

	var estimated struct {
		Gas uint64 `json:"gas"`
	}
	if err := json.NewDecoder(estResp.Body).Decode(&estimated); err != nil {
		t.Fatalf("decode estimate response: %v", err)
	}
	if estimated.Gas != 21000 {
		t.Fatalf("gas = %d, want 21000", estimated.Gas)
	}
}

func TestGateway_RevokedAPIKey_LosesAccess(t *testing.T) {
	srv, apiKeyService := newTestGateway(t)
	secret := issueTestAPIKey(t, srv)

	// Confirm access works before revoking.
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/v1/testnet/wallets/0xabc/balance", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("get balance: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status before revoke = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	key, err := apiKeyService.Authenticate(context.Background(), secret)
	if err != nil {
		t.Fatalf("authenticate to find key id: %v", err)
	}
	if err := apiKeyService.Revoke(context.Background(), key.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	req2, _ := http.NewRequest(http.MethodGet, srv.URL+"/v1/testnet/wallets/0xabc/balance", nil)
	req2.Header.Set("Authorization", "Bearer "+secret)
	resp2, err := srv.Client().Do(req2)
	if err != nil {
		t.Fatalf("get balance after revoke: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status after revoke = %d, want %d", resp2.StatusCode, http.StatusUnauthorized)
	}
}
