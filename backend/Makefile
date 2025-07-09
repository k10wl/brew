BUILD_DIR=bin
CMD_DIR=cmd
HTTP_ENTRYPOINT=$(CMD_DIR)/http.go
HTTP_BINARY=brew-http

# Go variables
GO=go
LDFLAGS=-ldflags="-s -w"

.PHONY: build-http run-http dev-http test coverage clean

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

build-http: $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(HTTP_BINARY) $(HTTP_ENTRYPOINT)

run-http: build-http
	./$(BUILD_DIR)/$(HTTP_BINARY)

dev-http:
	air -c .air.http.toml

test:
	$(GO) test ./...

coverage:
	@mkdir -p tmp
	$(GO) test -coverprofile=tmp/coverage.out ./...
	$(GO) tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated: tmp/coverage.html"

clean:
	rm -rf tmp/*
