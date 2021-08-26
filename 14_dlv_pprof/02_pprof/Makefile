.PHONY: run
run:
	go run .

.PHONY: ab
ab:
	@docker run \
		--net=host \
		--rm \
		-v `pwd`/request.json:/etc/request.json \
		jordi/ab -l -p /etc/request.json -s 60s -k -c 30 -n 10000 http://127.0.0.1:8080/mining

.PHONY: pprof-goroutine
pprof-goroutine:
	go tool pprof http://127.0.0.1:8080/debug/pprof/goroutine

.PHONY: pprof-heap
pprof-heap:
	go tool pprof -inuse_space http://127.0.0.1:8080/debug/pprof/heap

.PHONY: pprof-web
pprof-web:
	go tool pprof -http=:8081 http://127.0.0.1:8080/debug/pprof/profile