// Package arc implements the chainport.Adapter for the ARC Network,
// Aureon's phase-one blockchain.
package arc

import (
	"context"
	"fmt"

	chainport "github.com/jeielsantos/aureon/modules/contracts"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Network() chainport.Network {
	return "arc"
}

func (a *Adapter) CreateWallet(ctx context.Context) (chainport.Address, error) {
	return "", fmt.Errorf("arc: CreateWallet not implemented")
}

func (a *Adapter) GetBalance(ctx context.Context, address chainport.Address) (chainport.Balance, error) {
	return chainport.Balance{}, fmt.Errorf("arc: GetBalance not implemented")
}

func (a *Adapter) SendTransaction(ctx context.Context, tx chainport.Transaction) (chainport.TxHash, error) {
	return "", fmt.Errorf("arc: SendTransaction not implemented")
}

func (a *Adapter) EstimateGas(ctx context.Context, tx chainport.Transaction) (uint64, error) {
	return 0, fmt.Errorf("arc: EstimateGas not implemented")
}
