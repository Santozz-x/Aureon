package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Santozz-x/Aureon/modules/platform/internal/domain"
	"github.com/Santozz-x/Aureon/modules/platform/internal/usecase"
)

type APIKeyHandler struct {
	service *usecase.APIKeyService
}

func NewAPIKeyHandler(service *usecase.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: service}
}

type issueAPIKeyRequest struct {
	AccountID string `json:"account_id"`
}

type issueAPIKeyResponse struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

func (h *APIKeyHandler) Issue(w http.ResponseWriter, r *http.Request) {
	var req issueAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountID == "" {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	id, secret, err := h.service.Issue(r.Context(), domain.AccountID(req.AccountID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, issueAPIKeyResponse{ID: string(id), Secret: secret})
}

func (h *APIKeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := domain.APIKeyID(r.PathValue("id"))
	if err := h.service.Revoke(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
