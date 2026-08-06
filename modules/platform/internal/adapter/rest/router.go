package rest

import "net/http"

func NewRouter(walletHandler *WalletHandler, transactionHandler *TransactionHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /v1/{network}/wallets", walletHandler.CreateWallet)
	mux.HandleFunc("GET /v1/{network}/wallets/{address}/balance", walletHandler.GetBalance)

	mux.HandleFunc("POST /v1/{network}/transactions", transactionHandler.Send)
	mux.HandleFunc("POST /v1/{network}/transactions/estimate", transactionHandler.EstimateGas)

	return mux
}
