.PHONY: test lint fmt

test: lint fmt
	go test -v ./...

lint:
	go vet ./...
	staticcheck ./...

fmt:
	gofmt -l -w .
