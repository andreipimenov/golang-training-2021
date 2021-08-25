.PHONY: run
run:
	@export EXTERNAL_API_TOKEN=.token && export DB_CONN_STRING=.db_conn && export SECRET=.secret && go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: docker-build
docker-build:
	@docker build -t stock-service -f Dockerfile .

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

.PHONY: gen-mocks
gen-mocks:
	@docker run -v `pwd`:/src -w /src vektra/mockery:v2.7 --case snake --dir internal --output internal/mock --outpkg mock --all

.PHONY: run-db
run-db:
	@docker run \
		-d \
		-v `pwd`/db:/docker-entrypoint-initdb.d/ \
		--rm \
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=backend \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:12

.PHONY: docker-build-generator
docker-build-generator:
	@docker build -t generator -f Dockerfile.generator .

.PHONE: gen-proto
gen-proto: docker-build-generator
	@docker run -d --rm \
		-v `pwd`/api:/api \
		-v `pwd`/internal/pb:/pb \
		generator protoc \
		--go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true \
		-I /usr/local/include/. \
    	-I /api/. api.proto
