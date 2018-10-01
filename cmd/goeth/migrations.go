package main

import (
	"github.com/dneprix/goeth/model"
	"github.com/jinzhu/gorm"
)

func migrations(db *gorm.DB) {
	if *dbReset {
		db.DropTableIfExists(
			&model.Account{},
			&model.Transaction{},
		)
	}
	db.AutoMigrate(
		&model.Account{},
		&model.Transaction{},
	)

	return
}
