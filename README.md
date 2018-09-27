## GoETH

Service for managing ETH wallets

### Install and run

Create postgreSQL database (default `goeth`)

```
go get -u github.com/golang/dep
go get -u github.com/dneprix/goeth
cd $GOPATH/go/src/github.com/dneprix/goeth

make install-deps
make install

goeth --rpc-eth-node="127.0.0.1:8545"

```

### Run ethereum node (if you need it)

```
geth --testnet --rpc --syncmode="light"
```

### Create ethereum accounts (if you don't have)

```
geth --testnet account new
```

### Testing service
```
make test
```
