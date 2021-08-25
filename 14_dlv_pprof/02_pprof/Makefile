.PHONY: run
run:
	go run .

.PHONY: ab
ab:
	@docker run \
		--net=host \
		--rm \
		-v `pwd`/request.json:/etc/request.json \
		jordi/ab -p /etc/request.json -k -c 30 -n 10000 http://host.docker.internal:8080/mining

.PHONY: pprof-goroutine
pprof-goroutine:
	go tool pprof http://127.0.0.1:8080/debug/pprof/goroutine

.PHONY: pprof-heap
pprof-heap:
	go tool pprof -inuse_space http://127.0.0.1:8080/debug/pprof/heap
