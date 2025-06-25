lint:
	@echo "  > linting..."
	go vet ./...
	staticcheck ./...
	golangci-lint run ./...
	@echo "  > linting successfully finished"

test:
	@echo "  > testing..."
	go test -gcflags="-l" -race -v ./...
	@echo "  > testing successfully finished"

all:	
	make lint
	make test

.PHONY: lint test all