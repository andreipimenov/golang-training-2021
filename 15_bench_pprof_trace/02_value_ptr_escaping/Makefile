.PHONY: bench
bench:
	go test \
		-bench . \
		-benchmem \
		-trace t.out

.PHONY: escape
escape:
	go build -gcflags "-m" user.go

.PHONY: trace
trace:
	go tool trace t.out

.PHONY: clean
clean:
	rm t.out *.test