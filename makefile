APP_NAME=sylmark
GO_CMD=go
FRONTEND_DIR=server/sylgraph
BUILD_DIR=build

# Detect GOPATH automatically if not set
GOPATH := $(shell $(GO_CMD) env GOPATH)

.PHONY: all frontend backend clean run move install

all: frontend backend

# Build frontend
frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build

# Build Go backend
backend:
	mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(BUILD_DIR)/$(APP_NAME) .

# Install binary into GOPATH/bin
install: frontend backend
	mkdir -p $(GOPATH)/bin
	cp $(BUILD_DIR)/$(APP_NAME) $(GOPATH)/bin/$(APP_NAME)


# Cleanup
clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist
