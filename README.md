# GoETH

GoETH - service for managing ETH wallets.
1. Service communicates with ethereum geth node and uses IPC socket. Also service gets new transactions throw block subscription. Please set `--ipc-path=[$datadir]/geth.ipc` or `-p [$datadir]/geth.ipc`
2. Service download history of income transactions from last blocks (default 100). Please set `--load-history=[number]` or `-l [number]`

## Install, test and run

1. Create `postgreSQL` database (default `goeth`)
2. Run ethereum node `geth`
3. Install and run goETH service
```
$ go get -u github.com/golang/dep
$ go get -u github.com/dneprix/goeth
$ cd $GOPATH/go/src/github.com/dneprix/goeth

$ make install-deps
$ make install
$ make test

$ goeth --ipc-path=[$datadir]/geth.ipc
```

## Configuration

```
goeth -h

Usage: goeth [OPTIONS]
A service for handling ETH wallets

Options:              
   --public-addr   Listen address for API (env $GOETH_PUBLIC_ADDR) (default "0.0.0.0:8090")
-i --ipc-path      IPC socket path for ethereum geth node (env $GOETH_IPC_PATH)
-l --load-history  Number blocks in the past for loading transaction history (env $GOETH_LOAD_HISTORY) (default 100 blocks)
   --db-host       DB server host (env $GOETH_DB_HOST) (default "/run/postgresql")
   --db-port       DB server port (env $GOETH_DB_PORT)
   --db-user       DB server user (env $GOETH_DB_USER)
   --db-password   DB server password (env $GOETH_DB_PASSWORD)
   --db-name       Server DB name (env $GOETH_DB_NAME) (default "goeth")
   --db-reset      Remove db data before start service (env $GOETH_DB_RESET)

```

## API

```
POST /api/v1/SendEth
GET /api/v1/GetLast
```

### POST /api/v1/SendEth
Send ETH from one address to another address.
*Required fields:*
* `from` - ETH Address ("0x...")
* `to` - ETH Address ("0x...")
* `amount` - amount ETH (float64)

*Example request:*
```json
curl -X POST --data '{
  "from": "0xa70ee95f624120cf53d97d46c8063d077282de93",
  "to": "0xdaf39b5dcaa4f8dfbb07909d78d7e530e5b93484",
  "amount": 0.01
}' http://localhost:8090/api/v1/SendEth
```
*Example response:*
```json
{
  "tx_hash":"0x01a534848c561b68226d9234d84d7104bfa3503861fffb2a65a85621c2f487b0"
}
```

### GET /api/v1/GetLast
Show last income payments. `amount` - Ether value
*Example response:*
```json
[
  {
    "date": "Wed Oct 10 11:33:20 EEST 2018",
    "address": "0xdaf39b5dcaa4f8dfbb07909d78d7e530e5b93484",
    "amount": 0.01,
    "confirmations": "3922"
  },
  {
    "date": "Wed Oct 10 10:20:40 EEST 2018",
    "address": "0xdaf39b5dcaa4f8dfbb07909d78d7e530e5b93484",
    "amount": 0.024,
    "confirmations": "3890"
  }
]
```

# Performance testing
`
$ ab -c 100 -n 1000 http://localhost:8090/api/v1/GetLast
Concurrency Level:      100
Time taken for tests:   0.594 seconds
Complete requests:      1000
Failed requests:        0
`

`
$ ab -p files/test.txt -T application/json -c 1000 -n 1000 http://localhost:8090/api/v1/SendEth

Concurrency Level:      1000
Time taken for tests:   0.946 seconds
Complete requests:      1000
Failed requests:        0
`

# Notes
1. I didn't create Dockerfile and image to save my time. But I can create if you need it.
2. I develop service that load history process finished before starting http server. If history is not needed on the start I can put load history process in background to goroutine.
