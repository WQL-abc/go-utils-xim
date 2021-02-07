.PHONY: test
test:
	go test ./... -v

.PHONY: lint
lint:
	golangci-lint run --config ".github/.golangci.yml" --fast
