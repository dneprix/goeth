package model

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

const TX_CONFIRM_MAX = 6
const TX_CONFIRM_MIN = 3

type Transaction struct {
	Hash          common.Hash     `gorm:"type:bytea;primary_key"`
	AccountFrom   common.Address  `gorm:"type:bytea;"`
	AccountTo     common.Address  `gorm:"type:bytea;index:idx_account_to_view;"`
	Amount        *Wei            `gorm:"type:numeric;"`
	BlockNum      decimal.Decimal `gorm:"type:numeric"`
	Confirmations decimal.Decimal `gorm:"type:numeric"`
	ChainId       decimal.Decimal `gorm:"type:numeric"`
	Gas           uint64          `gorm:"type:bigint"`
	GasPrice      decimal.Decimal `gorm:"type:numeric"`
	Nonce         uint64          `gorm:"type:bigint"`
	Protected     bool            `gorm:"type:boolean"`
	View          bool            `gorm:"type:boolean;index:idx_account_to_view;"`
	CreatedAt     time.Time
}
