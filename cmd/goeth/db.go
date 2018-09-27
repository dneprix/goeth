package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func dbNew(host, port, user, password, name string) (*gorm.DB, error) {

	connect := "host=" + host + " dbname=" + name
	if len(port) > 0 {
		connect += " port=" + port
	}
	if len(user) > 0 {
		connect += " user=" + user
	}
	if len(password) > 0 {
		connect += " password=" + password
	}

	db, err := gorm.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	return db, nil
}
