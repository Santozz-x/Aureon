package usecase

import (
	"context"

	chainport "github.com/jeielsantos/aureon/modules/contracts"
)

type WalletService struct {
	adapterRegistry
}

func NewWalletService(adapters map[chainport.Network]chainport.Adapter) *WalletService {
	return &WalletService{adapterRegistry: newAdapterRegistry(adapters)}
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
