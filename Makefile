.PHONY: run
run:
	@go run cmd/main.go

.PHONY: run/env
run/env:
	@export EXTERNAL_API_TOKEN=.token && export DB_CONN_STRING=.db_conn && go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/main.go

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: docker/build
docker/build:
	@docker build -t stock-service .

.PHONY: docker-run
docker-run: docker-rm
	@docker run \
		--name stock-service \
		-p 80:8080 \
		-v `pwd`/secret:/secret \
		-e EXTERNAL_API_TOKEN=/secret/.token \
		-e DB_CONN_STRING=/secret/.db \
		stock-service

.PHONY: docker-stop
docker-stop:
	@docker stop stock-service || true

.PHONY: docker-rm
docker-rm: docker-stop
	@docker rm stock-service || true

.PHONY: gen-mocks
gen-mocks:
	@docker run \
	--name mockery \
	--rm \
	-v `pwd`:/src \
	-w /src \
	vektra/mockery:v2.7 --case snake --dir internal --output internal/mock --outpkg mock --all

.PHONY: run-db
run-db:
	@docker run \
		-d \
		-v `pwd`/db:/docker-entrypoint-initdb.d/ \
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=backend \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:12

.PHONY: mongo/run
mongo/run: mongo/rm
	@docker run \
		-d \
		-p 27017:27017 \
		--name mongo \
		-e MONGO_INITDB_ROOT_USERNAME=root \
        -e MONGO_INITDB_ROOT_PASSWORD=example \
	    mongo:5.0.2

.PHONY: mongo/rm
mongo/rm: mongo/stop
	@docker rm mongo || true

.PHONY: mongo/stop
mongo/stop:
	@docker stop mongo || true


.PHONY: docker/stack/deploy
docker/stack/deploy: docker/stack/rm docker/build
	@docker stack deploy -c docker-stack.yaml xxx || true

.PHONY: docker/stack/rm
docker/stack/rm:
	@docker stack rm xxx || true
