BUILD_DIR=bin
CMD_DIR=cmd
HTTP_ENTRYPOINT=$(CMD_DIR)/http.go
HTTP_BINARY=brew-http

# Go variables
GO=go
LDFLAGS=-ldflags="-s -w"

.PHONY: build-http run-http

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

build-http: $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(HTTP_BINARY) $(HTTP_ENTRYPOINT)

run-http: build-http
	./$(BUILD_DIR)/$(HTTP_BINARY)
