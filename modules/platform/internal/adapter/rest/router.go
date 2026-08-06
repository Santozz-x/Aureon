package rest

import "net/http"

// protect wraps a handler with authentication (usually middleware.RequireAPIKey).
func NewRouter(walletHandler *WalletHandler, transactionHandler *TransactionHandler, apiKeyHandler *APIKeyHandler, protect func(http.Handler) http.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Issuing/revoking API keys is intentionally not behind RequireAPIKey
	// itself — there is no bootstrap flow yet to obtain the very first
	// key. This must move behind a real Identity/signup flow before any
	// production use; see docs/tradeoffs.md.
	mux.HandleFunc("POST /v1/apikeys", apiKeyHandler.Issue)
	mux.HandleFunc("DELETE /v1/apikeys/{id}", apiKeyHandler.Revoke)

	mux.Handle("POST /v1/{network}/wallets", protect(http.HandlerFunc(walletHandler.CreateWallet)))
	mux.Handle("GET /v1/{network}/wallets/{address}/balance", protect(http.HandlerFunc(walletHandler.GetBalance)))
	mux.Handle("POST /v1/{network}/transactions", protect(http.HandlerFunc(transactionHandler.Send)))
	mux.Handle("POST /v1/{network}/transactions/estimate", protect(http.HandlerFunc(transactionHandler.EstimateGas)))

	return mux
}
