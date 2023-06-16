.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test-unit
test-unit:
	go test ./pkg/...

.PHONY: test-integration
test-integration:
	go test ./tests/...

.PHONY: test
test: test-unit test-integration
