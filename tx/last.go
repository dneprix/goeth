package tx

import (
	"github.com/dneprix/goeth/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
)

type Last interface {
	View() []*model.Transaction
	Sync(tx *model.Transaction)
	Load(accounts []common.Address)
}

type last struct {
	db       *gorm.DB
	list     []*model.Transaction
	accounts []common.Address
}

func NewLast(db *gorm.DB) Last {
	return &last{
		db: db,
	}
}

func (l *last) View() []*model.Transaction {

	return l.list
}

func (l *last) Sync(tx *model.Transaction) {

	return
}

func (l *last) Load(accounts []common.Address) {
	l.accounts = accounts
	l.db.Where(
		"account_to IN (?) AND (confirmations < ? OR view = ?)",
		l.accounts,
		model.TX_CONFIRM_MIN,
		false,
	).Find(&l.list)
}
