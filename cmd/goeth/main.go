package main

import (
	"fmt"
	"os"

	"github.com/dneprix/goeth"
	"github.com/gin-gonic/gin"
	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
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
		err = fmt.Errorf("failed to connect DB: %v", err)
		log.Fatalln(err)
	}
	migrations(db)

	client, err := goeth.NewClient(*ipcPath)
	if err != nil {
		log.Fatalln("failed to connect IPC socket:", err)
	}

	svc, err := goeth.NewService(db, client)
	if err != nil {
		log.Fatalln("failed to init service:", err)
	}

	err = svc.LoadAccounts()
	if err != nil {
		log.Fatalln("failed to add accounts:", err)
	}

	if *historyBlocks > 0 {
		log.Infof("load history form last %v blocks. Please wait", *historyBlocks)
		svc.LoadHistoryBlocks(*historyBlocks)
	}
	svc.LoadLast()
	go svc.LoadBalances()
	go svc.LoadNewBlocks()

	r := gin.Default()
	setHandlers(r, svc)
	if err := r.Run(*publicAddr); err != nil {
		log.Fatalln("HTTP server error:", err)
	}
}
