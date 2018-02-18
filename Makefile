.PHONY: list test


lint:
	gofmt -w cmd
	goimports -w cmd
	go fix cmd
	go vet

test:
	go test ./...
