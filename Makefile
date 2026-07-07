.PHONY: build test run lint clean

build:
	go build -o taskrunner ./cmd/taskrunner

test:
	go test ./... -v

run: build
	./taskrunner -file tasks.json -workers 3 -verbose

lint:
	go vet ./...
	@echo "--- gofmt ---"
	@diff=$$(gofmt -d .); \
	if [ -n "$$diff" ]; then \
		echo "$$diff"; \
		exit 1; \
	fi
	@echo "OK"

clean:
	rm -f taskrunner METRICS.md
