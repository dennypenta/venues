.PHONY: list test push


lint:
	gofmt -w cmd
	goimports -w cmd
	go fix cmd
	go vet
