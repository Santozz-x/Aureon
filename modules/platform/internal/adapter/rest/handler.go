package rest

import (
	"encoding/json"
	"net/http"

	chainport "github.com/Santozz-x/Aureon/modules/contracts"
	"github.com/Santozz-x/Aureon/modules/platform/internal/usecase"
)

type WalletHandler struct {
	service *usecase.WalletService
}

func NewWalletHandler(service *usecase.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

type createWalletResponse struct {
	Address string `json:"address"`
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	network := chainport.Network(r.PathValue("network"))

	address, err := h.service.CreateWallet(r.Context(), network)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, createWalletResponse{Address: string(address)})
}

type balanceResponse struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	network := chainport.Network(r.PathValue("network"))
	address := chainport.Address(r.PathValue("address"))

	balance, err := h.service.GetBalance(r.Context(), network, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, balanceResponse{Amount: balance.Amount, Symbol: balance.Symbol})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Headers are already sent; nothing meaningful to do with an encode
	// error here (e.g. client disconnected) beyond not silently ignoring it.
	_ = json.NewEncoder(w).Encode(payload)
}
