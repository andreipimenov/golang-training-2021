.PHONY: run
run:
	@export EXTERNAL_API_TOKEN=.token && go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: test
test:
	go test -v ./...