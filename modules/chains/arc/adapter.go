// Package arc implements the chainport.Adapter for the ARC Network,
// Aureon's phase-one blockchain.
package arc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	addr, err := toAddress(address)
	if err != nil {
		return chainport.Balance{}, err
	}

	balance, err := a.rpcClient.BalanceAt(ctx, addr)
	if err != nil {
		return chainport.Balance{}, fmt.Errorf("arc: get balance for %s: %w", address, err)
	}

	// Amount is the raw smallest-unit integer returned by eth_getBalance.
	// Whether Arc's native (USDC-denominated) balance uses 18 or 6 decimals
	// is not yet confirmed — see the follow-up in docs/networks/arc.md.
	// Callers must not assume a decimal scale until that's resolved.
	return chainport.Balance{Amount: balance.String(), Symbol: "USDC"}, nil
}

func (a *Adapter) EstimateGas(ctx context.Context, tx chainport.Transaction) (uint64, error) {
	msg, err := a.callMsg(tx)
	if err != nil {
		return 0, err
	}

	gas, err := a.rpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("arc: estimate gas from %s to %s: %w", tx.From, tx.To, err)
	}
	return gas, nil
}

func (a *Adapter) SendTransaction(ctx context.Context, tx chainport.Transaction) (chainport.TxHash, error) {
	msg, err := a.callMsg(tx)
	if err != nil {
		return "", err
	}

	privateKeyBytes, err := a.keyStore.Get(ctx, tx.From)
	if err != nil {
		return "", fmt.Errorf("arc: no key stored for %s: %w", tx.From, err)
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("arc: stored key for %s is invalid: %w", tx.From, err)
	}

	nonce, err := a.rpcClient.PendingNonceAt(ctx, msg.From)
	if err != nil {
		return "", fmt.Errorf("arc: get nonce for %s: %w", tx.From, err)
	}

	gasPrice, err := a.rpcClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("arc: suggest gas price: %w", err)
	}

	gasLimit, err := a.rpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("arc: estimate gas from %s to %s: %w", tx.From, tx.To, err)
	}

	chainID, err := a.rpcClient.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("arc: chain id: %w", err)
	}

	unsignedTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       msg.To,
		Value:    msg.Value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     msg.Data,
	})

	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("arc: sign transaction from %s: %w", tx.From, err)
	}

	if err := a.rpcClient.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("arc: send transaction from %s: %w", tx.From, err)
	}

	return chainport.TxHash(signedTx.Hash().Hex()), nil
}

// callMsg validates and converts a chainport.Transaction into the
// go-ethereum call message shared by EstimateGas and SendTransaction.
func (a *Adapter) callMsg(tx chainport.Transaction) (ethereum.CallMsg, error) {
	from, err := toAddress(tx.From)
	if err != nil {
		return ethereum.CallMsg{}, err
	}
	to, err := toAddress(tx.To)
	if err != nil {
		return ethereum.CallMsg{}, err
	}

	value := big.NewInt(0)
	if tx.Value != "" {
		v, ok := new(big.Int).SetString(tx.Value, 10)
		if !ok {
			return ethereum.CallMsg{}, fmt.Errorf("arc: invalid transaction value %q", tx.Value)
		}
		value = v
	}

	return ethereum.CallMsg{From: from, To: &to, Value: value, Data: tx.Data}, nil
}

func toAddress(address chainport.Address) (common.Address, error) {
	if !common.IsHexAddress(string(address)) {
		return common.Address{}, fmt.Errorf("arc: invalid address %q", address)
	}
	return common.HexToAddress(string(address)), nil
}
