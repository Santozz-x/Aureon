// Package rpc wraps the Arc network's JSON-RPC endpoint, which follows the
// standard Ethereum eth_* API (see docs/networks/arc.md at the repo root).
package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	eth *ethclient.Client
}

func Dial(ctx context.Context, rpcURL string) (*Client, error) {
	eth, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("rpc: dial %s: %w", rpcURL, err)
	}
	return &Client{eth: eth}, nil
}

func (c *Client) Close() {
	c.eth.Close()
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	id, err := c.eth.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("rpc: chain id: %w", err)
	}
	return id, nil
}

func (c *Client) BalanceAt(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.eth.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("rpc: balance at %s: %w", address, err)
	}
	return balance, nil
}

func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	gas, err := c.eth.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("rpc: estimate gas: %w", err)
	}
	return gas, nil
}

func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := c.eth.SendTransaction(ctx, tx); err != nil {
		return fmt.Errorf("rpc: send transaction: %w", err)
	}
	return nil
}
