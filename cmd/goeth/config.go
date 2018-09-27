package main

import (
	cli "github.com/jawher/mow.cli"
)

// Service goETH config
var (
	publicAddr = app.String(cli.StringOpt{
		Name:   "public-addr",
		Desc:   "Listen address for external access and public HTTP API",
		EnvVar: "GOETH_PUBLIC_ADDR",
		Value:  "0.0.0.0:8090",
	})
	rpcEthNode = app.String(cli.StringOpt{
		Name:   "rpc-eth-node",
		Desc:   "RPC ethereum node address",
		EnvVar: "GOETH_RPC_ETH_NODE",
		Value:  "http://0.0.0.0:8545",
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
)
