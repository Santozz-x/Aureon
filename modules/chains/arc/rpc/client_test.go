package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type rpcRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

func newMockServer(t *testing.T, handler func(method string) (interface{}, *rpcError)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		result, rpcErr := handler(req.Method)
		resp := rpcResponse{JSONRPC: "2.0", ID: req.ID}
		if rpcErr != nil {
			resp.Error = rpcErr
		} else {
			resp.Result = result
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
}

func TestClient_ChainID(t *testing.T) {
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		if method != "eth_chainId" {
			return nil, &rpcError{Code: -32601, Message: "method not found: " + method}
		}
		return fmt.Sprintf("0x%x", 5042002), nil
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	id, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatalf("ChainID: %v", err)
	}
	if id.Cmp(big.NewInt(5042002)) != 0 {
		t.Fatalf("ChainID = %s, want 5042002", id)
	}
}

func TestClient_BalanceAt(t *testing.T) {
	want := big.NewInt(1_000_000_000_000_000_000)
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		if method != "eth_getBalance" {
			return nil, &rpcError{Code: -32601, Message: "method not found: " + method}
		}
		return fmt.Sprintf("0x%x", want), nil
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	address := common.HexToAddress("0x00000000000000000000000000000000000AAA")
	got, err := client.BalanceAt(context.Background(), address)
	if err != nil {
		t.Fatalf("BalanceAt: %v", err)
	}
	if got.Cmp(want) != 0 {
		t.Fatalf("BalanceAt = %s, want %s", got, want)
	}
}

func TestClient_BalanceAt_UpstreamError(t *testing.T) {
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		return nil, &rpcError{Code: -32000, Message: "header not found"}
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	_, err = client.BalanceAt(context.Background(), common.HexToAddress("0xAA"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_EstimateGas(t *testing.T) {
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		if method != "eth_estimateGas" {
			return nil, &rpcError{Code: -32601, Message: "method not found: " + method}
		}
		return fmt.Sprintf("0x%x", 21000), nil
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	from := common.HexToAddress("0x000000000000000000000000000000000000A0")
	to := common.HexToAddress("0x000000000000000000000000000000000000B0")
	gas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: from, To: &to})
	if err != nil {
		t.Fatalf("EstimateGas: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("EstimateGas = %d, want 21000", gas)
	}
}

func TestClient_EstimateGas_UpstreamError(t *testing.T) {
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		return nil, &rpcError{Code: -32000, Message: "execution reverted"}
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	from := common.HexToAddress("0x000000000000000000000000000000000000A0")
	to := common.HexToAddress("0x000000000000000000000000000000000000B0")
	_, err = client.EstimateGas(context.Background(), ethereum.CallMsg{From: from, To: &to})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_PendingNonceAt(t *testing.T) {
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		if method != "eth_getTransactionCount" {
			return nil, &rpcError{Code: -32601, Message: "method not found: " + method}
		}
		return fmt.Sprintf("0x%x", 7), nil
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress("0x000000000000000000000000000000000000A0"))
	if err != nil {
		t.Fatalf("PendingNonceAt: %v", err)
	}
	if nonce != 7 {
		t.Fatalf("PendingNonceAt = %d, want 7", nonce)
	}
}

func TestClient_SuggestGasPrice(t *testing.T) {
	want := big.NewInt(1_000_000_000)
	srv := newMockServer(t, func(method string) (interface{}, *rpcError) {
		if method != "eth_gasPrice" {
			return nil, &rpcError{Code: -32601, Message: "method not found: " + method}
		}
		return fmt.Sprintf("0x%x", want), nil
	})
	defer srv.Close()

	client, err := Dial(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("SuggestGasPrice: %v", err)
	}
	if price.Cmp(want) != 0 {
		t.Fatalf("SuggestGasPrice = %s, want %s", price, want)
	}
}
