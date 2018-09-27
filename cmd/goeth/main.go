package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dneprix/goeth"
	"github.com/dneprix/goeth/model"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/jawher/mow.cli"
)

var app = cli.App("goeth", "A service for handling ETH wallets")

func main() {
	app.Action = appMain
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func appMain() {
	db, err := dbNew(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	if err != nil {
		err = fmt.Errorf("[ERR] failed to connect DB: %v", err)
		log.Fatalln(err)
	}
	db.AutoMigrate(&model.Wallet{}, &model.Transactions{})

	client, err := ethclient.Dial(*rpcEthNode)
	if err != nil {
		log.Fatalln("[ERR] failed to connect RPC node:", err)
	}

	svc, err := goeth.NewService(db, client)
	if err != nil {
		log.Fatalln("[ERR] failed to init service:", err)
	}

	svc.SyncWallets()

	r := gin.Default()
	SetHandlers(r, svc)
	if err := r.Run(*publicAddr); err != nil {
		log.Fatalln("[ERR] HTTP server:", err)
	}
}
