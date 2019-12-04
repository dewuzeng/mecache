.PHONY: fmt test

GOFILES=$(shell find . -type f -name '*.go' | grep -v /vendor/)

fmt:
	@goimports -l -w ${GOFILES}
	@gofmt -l -w -s ${GOFILES}
	@go vet -v $(go list ./...| grep -v /vendor/)

test:
	go test -count=1 -v ./...