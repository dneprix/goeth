package tx

import (
	"github.com/dneprix/goeth/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
)

type Last interface {
	View() []*model.Transaction
	Load(accounts []common.Address)
}

type last struct {
	db           *gorm.DB
	list         []*model.Transaction
	accounts     []common.Address
	saveChangesC chan *model.Transaction
}

func NewLast(db *gorm.DB) Last {
	return &last{
		db: db,
	}
}

func (l *last) View() []*model.Transaction {
	defer func() {
		for _, tx := range l.list {
			tx.View = true
		}
	}()
	return l.list
}

func (l *last) Load(accounts []common.Address) {
	l.accounts = accounts
	l.db.Where(
		"account_to IN (?) AND (confirmations < ? OR view = ?)",
		l.accounts,
		model.TX_CONFIRM_MIN,
		false,
	).Find(&l.list)
	go l.saveChanges()
}

func (l *last) saveChanges() {
	for {
		select {
		case tx, ok := <-l.saveChangesC:
			if !ok {
				return
			} else if tx == nil {
				continue
			}

			l.db.Save(tx)
		}
	}
	return
}
