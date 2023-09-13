build:
	go build -o bin/defi-poker

run: build
	@./bin/defi-poker

test:
	go test -v ./...

.PHONY: build