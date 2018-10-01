package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

type Wei string

func WeiFromString(value string) *Wei {
	w := (Wei)(value)
	return &w
}

func WeiFromEth(eth float64) *Wei {
	wei := decimal.NewFromFloat(eth).Mul(decimal.NewFromFloat(1e18)).String()
	return WeiFromString(wei)
}

func (w *Wei) String() string {
	return string(*w)
}

func (w *Wei) Ether() float64 {
	d, _ := decimal.NewFromString(w.String())
	eth, _ := d.Div(decimal.NewFromFloat(1e18)).Float64()
	return eth
}

func (w *Wei) Big() *hexutil.Big {
	bigint := big.NewInt(0)
	bigint.SetString(w.String(), 10)
	return (*hexutil.Big)(bigint)
}
