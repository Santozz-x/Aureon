package usecase

import (
	"context"
	"fmt"

	chainport "github.com/jeielsantos/aureon/modules/contracts"
)

type WalletService struct {
	adapters map[chainport.Network]chainport.Adapter
}

func NewWalletService(adapters map[chainport.Network]chainport.Adapter) *WalletService {
	return &WalletService{adapters: adapters}
}

func (s *WalletService) adapter(network chainport.Network) (chainport.Adapter, error) {
	adapter, ok := s.adapters[network]
	if !ok {
		return nil, fmt.Errorf("usecase: unsupported network %q", network)
	}
	return adapter, nil
}

func (s *WalletService) CreateWallet(ctx context.Context, network chainport.Network) (chainport.Address, error) {
	adapter, err := s.adapter(network)
	if err != nil {
		return "", err
	}
	return adapter.CreateWallet(ctx)
}

func (s *WalletService) GetBalance(ctx context.Context, network chainport.Network, address chainport.Address) (chainport.Balance, error) {
	adapter, err := s.adapter(network)
	if err != nil {
		return chainport.Balance{}, err
	}
	return adapter.GetBalance(ctx, address)
}
