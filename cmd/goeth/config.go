package main

import (
	cli "github.com/jawher/mow.cli"
)

// Service goETH config
var (
	publicAddr = app.String(cli.StringOpt{
		Name:   "public-addr",
		Desc:   "Listen address for API",
		EnvVar: "GOETH_PUBLIC_ADDR",
		Value:  "0.0.0.0:8090",
	})
	ipcPath = app.String(cli.StringOpt{
		Name:   "i ipc-path",
		Desc:   "IPC socket path for ethereum geth node",
		EnvVar: "GOETH_IPC_PATH",
		Value:  "",
	})
	historyBlocks = app.Int(cli.IntOpt{
		Name:   "l load-history",
		Desc:   "Number blocks in the past for loading transaction history",
		EnvVar: "GOETH_LOAD_HISTORY",
		Value:  100,
	})
)

// PostgreSql connection config
var (
	dbHost = app.String(cli.StringOpt{
		Name:   "db-host",
		Desc:   "DB server host",
		EnvVar: "GOETH_DB_HOST",
		Value:  "/run/postgresql", //localhost
	})
	dbPort = app.String(cli.StringOpt{
		Name:   "db-port",
		Desc:   "DB server port",
		EnvVar: "GOETH_DB_PORT",
		Value:  "", //5432
	})
	dbUser = app.String(cli.StringOpt{
		Name:   "db-user",
		Desc:   "DB server user",
		EnvVar: "GOETH_DB_USER",
		Value:  "",
	})
	dbPassword = app.String(cli.StringOpt{
		Name:   "db-password",
		Desc:   "DB server password",
		EnvVar: "GOETH_DB_PASSWORD",
		Value:  "",
	})
	dbName = app.String(cli.StringOpt{
		Name:   "db-name",
		Desc:   "Server DB name",
		EnvVar: "GOETH_DB_NAME",
		Value:  "goeth",
	})
	dbReset = app.Bool(cli.BoolOpt{
		Name:   "db-reset",
		Desc:   "Remove db data before start service",
		EnvVar: "GOETH_DB_RESET",
		Value:  false,
	})
)
