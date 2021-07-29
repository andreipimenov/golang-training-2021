.PHONY: run
run:
	@go run cmd/main.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: test
test:
	go test -v ./...