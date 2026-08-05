// Package chainport defines the blockchain-agnostic port that every
// network adapter (ARC Network, Ethereum, Solana, ...) must implement.
// It has no dependency on any concrete chain and is the contract shared
// between the platform module and every modules/chains/* module.
package chainport

import "context"

type Network string

type Address string

type TxHash string

type Balance struct {
	Amount string
	Symbol string
}

type Transaction struct {
	From  Address
	To    Address
	Value string
	Data  []byte
}

type Adapter interface {
	Network() Network
	CreateWallet(ctx context.Context) (Address, error)
	GetBalance(ctx context.Context, address Address) (Balance, error)
	SendTransaction(ctx context.Context, tx Transaction) (TxHash, error)
	EstimateGas(ctx context.Context, tx Transaction) (uint64, error)
}
