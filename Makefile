.PHONY: run
run:
	@EXTERNAL_API_TOKEN=./secret/.token \
		DB_CONN_STRING=./secret/.db_conn_localhost \
		go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: docker-build
docker-build:
	@docker build -t stock-service .

.PHONY: docker-run
docker-run:
	-docker network create stock-service-net
	@docker run \
		--name stock-service \
		-d \
		--rm \
		--net stock-service-net \
		-p 80:8080 \
		-v `pwd`/secret:/secret \
		-e EXTERNAL_API_TOKEN=/secret/.token \
		-e DB_CONN_STRING=/secret/.db_conn \
		stock-service

.PHONY: docker-stop
docker-stop:
	@-docker stop stock-service
	@-docker stop db
	@-docker network rm stock-service-net

.PHONY: gen-mocks
gen-mocks:
	@docker run -v `pwd`:/src -w /src vektra/mockery:v2.7 --case snake --dir internal --output internal/mock --outpkg mock --all

.PHONY: run-postgres
run-postgres:
	@-docker network create stock-service-net
	@docker run \
		-d \
		-v `pwd`/db:/docker-entrypoint-initdb.d/ \
		--rm \
		--net stock-service-net \
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=backend \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:12

.PHONY: run-mongo
run-mongo:
	@-docker network create stock-service-net
	@docker run \
		-d \
		-v `pwd`/db:/docker-entrypoint-initdb.d/ \
		--rm \
		--net stock-service-net \
		-p 27017:27017 \
		--name db \
		-e MONGO_INITDB_ROOT_USERNAME=root \
		-e MONGO_INITDB_ROOT_PASSWORD=password \
		mongo:5
