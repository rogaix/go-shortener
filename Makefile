# Variables
BINARY_NAME=url-shortener
BACKEND_DIR=./backend
MAIN_FILE=$(BACKEND_DIR)/cmd/server/main.go
PORT=8080

# Targets
.PHONY: all build clean test vet fmt run

all: build

build:
	@echo "Building backend..."
	@go build -o $(BACKEND_DIR)/$(BINARY_NAME) $(MAIN_FILE)

clean:
	@echo "Cleaning up..."
	@rm -f $(BACKEND_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./$(BACKEND_DIR)/...

vet:
	@echo "Running go vet..."
	@go vet ./$(BACKEND_DIR)/...

fmt:
	@echo "Formatting code..."
	@go fmt ./$(BACKEND_DIR)/...

run: build
	@echo "Starting backend server..."
	@$(BACKEND_DIR)/$(BINARY_NAME)
