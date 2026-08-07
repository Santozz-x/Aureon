package usecase

import (
	"context"

	chainport "github.com/Santozz-x/Aureon/modules/contracts"
)

type TransactionService struct {
	adapterRegistry
}

func NewTransactionService(adapters map[chainport.Network]chainport.Adapter) *TransactionService {
	return &TransactionService{adapterRegistry: newAdapterRegistry(adapters)}
}

func (s *TransactionService) Send(ctx context.Context, network chainport.Network, tx chainport.Transaction) (chainport.TxHash, error) {
	adapter, err := s.adapter(network)
	if err != nil {
		return "", err
	}
	return adapter.SendTransaction(ctx, tx)
}

func (s *TransactionService) EstimateGas(ctx context.Context, network chainport.Network, tx chainport.Transaction) (uint64, error) {
	adapter, err := s.adapter(network)
	if err != nil {
		return 0, err
	}
	return adapter.EstimateGas(ctx, tx)
}
