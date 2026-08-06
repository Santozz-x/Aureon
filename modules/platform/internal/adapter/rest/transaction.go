package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	chainport "github.com/jeielsantos/aureon/modules/contracts"
	"github.com/jeielsantos/aureon/modules/platform/internal/usecase"
)

type TransactionHandler struct {
	service *usecase.TransactionService
}

func NewTransactionHandler(service *usecase.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

type transactionRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data,omitempty"`
}

func (req transactionRequest) toDomain() (chainport.Transaction, error) {
	data, err := decodeHexData(req.Data)
	if err != nil {
		return chainport.Transaction{}, fmt.Errorf("invalid data: %w", err)
	}
	return chainport.Transaction{
		From:  chainport.Address(req.From),
		To:    chainport.Address(req.To),
		Value: req.Value,
		Data:  data,
	}, nil
}

type sendTransactionResponse struct {
	TxHash string `json:"tx_hash"`
}

func (h *TransactionHandler) Send(w http.ResponseWriter, r *http.Request) {
	network := chainport.Network(r.PathValue("network"))

	var req transactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := req.toDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	txHash, err := h.service.Send(r.Context(), network, tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, sendTransactionResponse{TxHash: string(txHash)})
}

type estimateGasResponse struct {
	Gas uint64 `json:"gas"`
}

func (h *TransactionHandler) EstimateGas(w http.ResponseWriter, r *http.Request) {
	network := chainport.Network(r.PathValue("network"))

	var req transactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := req.toDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gas, err := h.service.EstimateGas(r.Context(), network, tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, estimateGasResponse{Gas: gas})
}

// decodeHexData decodes an optional "0x"-prefixed hex string into bytes.
// An empty string is valid and yields nil data.
func decodeHexData(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		s = s[2:]
	}
	return hex.DecodeString(s)
}
