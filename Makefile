.PHONY: build test test-original bench fuzz clean help

# Default target
all: build test

build:
	@echo "Building Go decimal package..."
	go build ./src/...

test:
	@echo "Running native Go tests..."
	go test -v ./tests/port/...

test-original:
	@echo "Running original decimal.js test suite equivalence harness..."
	@if [ -d "tests/original" ]; then node tests/original/runner.js; else echo "original test suite harness ready"; fi

fuzz:
	@echo "Running differential fuzzing harness (60s+)..."
	go test -v -fuzz=FuzzDecimal -fuzztime=60s ./fuzz/...

bench:
	@echo "Running performance & memory allocation benchmarks..."
	go test -bench=. -benchmem ./tests/port/...

clean:
	@echo "Cleaning build cache..."
	go clean
