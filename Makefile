.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test-unit
test-unit:
	go test ./pkg/...

.PHONY: test
test: test-unit
