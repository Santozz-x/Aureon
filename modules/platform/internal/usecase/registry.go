package usecase

import (
	"fmt"

	chainport "github.com/Santozz-x/Aureon/modules/contracts"
)

type adapterRegistry struct {
	adapters map[chainport.Network]chainport.Adapter
}

func newAdapterRegistry(adapters map[chainport.Network]chainport.Adapter) adapterRegistry {
	return adapterRegistry{adapters: adapters}
}

func (r adapterRegistry) adapter(network chainport.Network) (chainport.Adapter, error) {
	adapter, ok := r.adapters[network]
	if !ok {
		return nil, fmt.Errorf("usecase: unsupported network %q", network)
	}
	return adapter, nil
}
