package proto

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type ApiSendEthRequest struct {
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Amount float64        `json:"amount"`
}

type ApiSendEthResponse struct {
	TxHash common.Hash `json:"tx_hash"`
}

type ApiGetLastItemResponse struct {
	Date          string          `json:"date"`
	Address       common.Address  `json:"address"`
	Amount        float64         `json:"amount"`
	Confirmations decimal.Decimal `json:"confirmations"`
}
