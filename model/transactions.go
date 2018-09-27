package model

import "github.com/jinzhu/gorm"

type Transactions struct {
	gorm.Model
	From     string
	To       string
	BlockNum int
	TxHash   string
	TxFee    int
}
