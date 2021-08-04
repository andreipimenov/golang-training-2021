.PHONY: run
run:
	@export EXTERNAL_API_TOKEN=.token && go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: test
test:
	go test -v ./...

.PHONY: docker-build
docker-build:
	@docker build -t stock-service .

.PHONY: docker-run
docker-run:
	@docker run \
		--name stock-service \
		-d \
		--rm \
		-p 80:8080 \
		-v `pwd`/secret:/secret \
		-e EXTERNAL_API_TOKEN=/secret/.token \
		stock-service

.PHONY: docker-stop
docker-stop:
	@docker stop stock-service