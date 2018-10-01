package model

import (
	"github.com/ethereum/go-ethereum/common"
)

type Account struct {
	Address common.Address `gorm:"primary_key;type:bytea;"`
	Balance *Wei           `gorm:"type:numeric;"`
}
