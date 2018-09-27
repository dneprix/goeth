install-deps:
	dep ensure

update-deps:
	dep ensure -update

install:
	go install github.com/dneprix/goeth/cmd/...

test:
	go test -v -race ./...
