// Package arc implements the chainport.Adapter for the ARC Network,
// Aureon's phase-one blockchain.
package arc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/jeielsantos/aureon/modules/chains/arc/rpc"
	chainport "github.com/jeielsantos/aureon/modules/contracts"
)

type Adapter struct {
	rpcClient *rpc.Client
	keyStore  chainport.KeyStore
}

func NewAdapter(rpcClient *rpc.Client, keyStore chainport.KeyStore) *Adapter {
	return &Adapter{rpcClient: rpcClient, keyStore: keyStore}
}

func (a *Adapter) Network() chainport.Network {
	return "arc"
}

func (a *Adapter) CreateWallet(ctx context.Context) (chainport.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("arc: generate key: %w", err)
	}

	address := chainport.Address(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	if err := a.keyStore.Put(ctx, address, crypto.FromECDSA(privateKey)); err != nil {
		return "", fmt.Errorf("arc: store key for %s: %w", address, err)
	}

	return address, nil
}

func (a *Adapter) GetBalance(ctx context.Context, address chainport.Address) (chainport.Balance, error) {
	if !common.IsHexAddress(string(address)) {
		return chainport.Balance{}, fmt.Errorf("arc: invalid address %q", address)
	}

	balance, err := a.rpcClient.BalanceAt(ctx, common.HexToAddress(string(address)))
	if err != nil {
		return chainport.Balance{}, fmt.Errorf("arc: get balance for %s: %w", address, err)
	}

	// Amount is the raw smallest-unit integer returned by eth_getBalance.
	// Whether Arc's native (USDC-denominated) balance uses 18 or 6 decimals
	// is not yet confirmed — see the follow-up in docs/networks/arc.md.
	// Callers must not assume a decimal scale until that's resolved.
	return chainport.Balance{Amount: balance.String(), Symbol: "USDC"}, nil
}

func (a *Adapter) SendTransaction(ctx context.Context, tx chainport.Transaction) (chainport.TxHash, error) {
	return "", fmt.Errorf("arc: SendTransaction not implemented")
}

func (a *Adapter) EstimateGas(ctx context.Context, tx chainport.Transaction) (uint64, error) {
	return 0, fmt.Errorf("arc: EstimateGas not implemented")
}
