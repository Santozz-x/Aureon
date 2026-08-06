package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_LogsMethodPathStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/v1/arc/wallets", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log line: %v (raw: %s)", err, buf.String())
	}

	if entry["method"] != http.MethodPost {
		t.Fatalf("method = %v, want %v", entry["method"], http.MethodPost)
	}
	if entry["path"] != "/v1/arc/wallets" {
		t.Fatalf("path = %v, want /v1/arc/wallets", entry["path"])
	}
	if entry["status"] != float64(http.StatusCreated) {
		t.Fatalf("status = %v, want %v", entry["status"], http.StatusCreated)
	}
	if entry["request_id"] == "" || entry["request_id"] == nil {
		t.Fatal("expected a non-empty request_id")
	}
}

func TestLogging_DefaultsStatusToOK(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handler never calls WriteHeader explicitly
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}
	if entry["status"] != float64(http.StatusOK) {
		t.Fatalf("status = %v, want %v", entry["status"], http.StatusOK)
	}
}

func TestLogging_PropagatesIncomingRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	var gotFromContext string
	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := RequestIDFromContext(r.Context())
		if !ok {
			t.Error("expected request id in context")
		}
		gotFromContext = id
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(requestIDHeader, "client-supplied-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if gotFromContext != "client-supplied-id" {
		t.Fatalf("context request id = %q, want %q", gotFromContext, "client-supplied-id")
	}
	if rec.Header().Get(requestIDHeader) != "client-supplied-id" {
		t.Fatalf("response header request id = %q, want %q", rec.Header().Get(requestIDHeader), "client-supplied-id")
	}
}

func TestLogging_GeneratesRequestIDWhenMissing(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get(requestIDHeader) == "" {
		t.Fatal("expected a generated request id in the response header")
	}
}
