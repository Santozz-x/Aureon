package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeielsantos/aureon/modules/platform/internal/infra/apikeystore"
	"github.com/jeielsantos/aureon/modules/platform/internal/usecase"
)

func TestRequireAPIKey_MissingHeader(t *testing.T) {
	service := usecase.NewAPIKeyService(apikeystore.NewMemory())
	handler := RequireAPIKey(service)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/arc/wallets/0x1/balance", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIKey_InvalidSecret(t *testing.T) {
	service := usecase.NewAPIKeyService(apikeystore.NewMemory())
	handler := RequireAPIKey(service)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/arc/wallets/0x1/balance", nil)
	req.Header.Set("Authorization", "Bearer aur_not-a-real-key")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIKey_ValidSecret(t *testing.T) {
	store := apikeystore.NewMemory()
	service := usecase.NewAPIKeyService(store)

	_, secret, err := service.Issue(context.Background(), "acct_1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	called := false
	handler := RequireAPIKey(service)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		key, ok := APIKeyFromContext(r.Context())
		if !ok {
			t.Fatal("expected api key in request context")
		}
		if key.AccountID != "acct_1" {
			t.Fatalf("context api key account = %s, want acct_1", key.AccountID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/arc/wallets/0x1/balance", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRequireAPIKey_RevokedSecret(t *testing.T) {
	store := apikeystore.NewMemory()
	service := usecase.NewAPIKeyService(store)

	id, secret, err := service.Issue(context.Background(), "acct_1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if err := service.Revoke(context.Background(), id); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	handler := RequireAPIKey(service)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/arc/wallets/0x1/balance", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
