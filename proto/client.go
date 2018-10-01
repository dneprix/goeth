package proto

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ClientTransaction struct {
	From        common.Address
	To          common.Address
	Value       *hexutil.Big
	BlockNumber *hexutil.Big
}
