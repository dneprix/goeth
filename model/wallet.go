package model

import "github.com/jinzhu/gorm"

type Wallet struct {
	gorm.Model
	Address string `gorm:"not null;unique"`
	Balance string
}
