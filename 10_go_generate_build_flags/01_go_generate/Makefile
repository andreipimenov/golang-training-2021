.PHONY: gen
gen:
	@TYPE=PriceService OUTPUT=../svc.go go generate gen/gen.go

.PHONY: run
run:
	@go run *.go

.PHONY: clean
clean:
	@rm svc.go