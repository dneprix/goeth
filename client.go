package goeth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dneprix/goeth/model"
	"github.com/dneprix/goeth/proto"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"
)

type Client interface {
	AccountsList() ([]common.Address, error)
	BalanceAt(addr common.Address) (*model.Wei, error)
	BlockByHash(hash common.Hash) (*types.Block, error)
	BlockByNumber(num *big.Int) (*types.Block, error)
	BlockNumByTxHash(txHash common.Hash) (blockNum *decimal.Decimal, err error)
	CheckPassphrase(addr common.Address, passphrase string) bool
	TransactionSender(tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	SendTransaction(from common.Address, to common.Address, amount *model.Wei, passphrase string) (*common.Hash, error)
	SubscribeNewHead(ch chan *types.Header) (ethereum.Subscription, error)
}

func NewClient(rawurl string) (Client, error) {
	rpcClient, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	client := &client{
		rpc: rpcClient,
		eth: ethclient.NewClient(rpcClient),
	}
	return client, nil
}

type client struct {
	rpc *rpc.Client
	eth *ethclient.Client
}

func (c *client) AccountsList() (accounts []common.Address, err error) {
	err = c.rpc.Call(&accounts, "eth_accounts", nil)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c *client) BalanceAt(addr common.Address) (*model.Wei, error) {
	balance, err := c.eth.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, err
	}
	return model.WeiFromString(balance.String()), nil
}

func (c *client) BlockByHash(hash common.Hash) (*types.Block, error) {
	return c.eth.BlockByHash(context.Background(), hash)
}

func (c *client) BlockByNumber(num *big.Int) (*types.Block, error) {
	return c.eth.BlockByNumber(context.Background(), num)
}

func (c *client) BlockNumByTxHash(txHash common.Hash) (blockNum *decimal.Decimal, err error) {
	var clientTx proto.ClientTransaction
	if err := c.rpc.Call(&clientTx, "eth_getTransactionByHash", txHash); err != nil {
		return nil, err
	}
	if clientTx.BlockNumber == nil {
		return nil, fmt.Errorf("no block number")
	}
	decBlockNumber := decimal.NewFromBigInt((*big.Int)(clientTx.BlockNumber), 0)
	return &decBlockNumber, nil
}

func (c *client) CheckPassphrase(addr common.Address, passphrase string) bool {
	if err := c.rpc.Call(nil, "personal_unlockAccount", addr, passphrase, 0); err != nil {
		return false
	}
	return true
}

func (c *client) TransactionSender(tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	return c.eth.TransactionSender(context.Background(), tx, block, index)
}

func (c *client) SendTransaction(from common.Address, to common.Address, amount *model.Wei, passphrase string) (*common.Hash, error) {
	var txHash common.Hash
	clientTx := proto.ClientTransaction{
		From:  from,
		To:    to,
		Value: amount.Big(),
	}
	if err := c.rpc.Call(&txHash, "eth_sendTransaction", clientTx); err != nil {
		return nil, err
	}
	return &txHash, nil
}

func (c *client) SubscribeNewHead(ch chan *types.Header) (ethereum.Subscription, error) {
	return c.eth.SubscribeNewHead(context.Background(), ch)
}
