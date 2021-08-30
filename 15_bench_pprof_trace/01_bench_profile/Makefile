.PHONY: bench
bench:
	go test \
		-bench . \
		-benchmem \
		-cpuprofile c.out \
		-memprofile m.out \

.PHONY: escape
escape:
	go build -gcflags "-m" block.go

.PHONY: cpu
cpu:
	go tool pprof c.out

.PHONY: mem
mem:
	go tool pprof m.out

.PHONY: cpu-http
cpu-http:
	go tool pprof -http :8080 c.out

.PHONY: mem-http
mem-http:
	go tool pprof -http :8080 m.out

.PHONY: clean
clean: 
	rm c.out m.out *.test