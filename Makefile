SHELL := /bin/bash

PRIVATE_KEY_FILE=ec_private_key.pem
PUBLIC_KEY_FILE=ec_public_key.pem
LINT_CONFIG_FILE=golangci.yml

all: generate-keys lint docker

docker:
	@echo "Building docker image..."
	docker-compose up
	@echo "Docker image built successfully"

generate-keys:
	@if [ ! -f $(PRIVATE_KEY_FILE) ]; then \
		echo "Generating ECDSA private key..."; \
		openssl ecparam -genkey -name prime256v1 -noout -out $(PRIVATE_KEY_FILE); \
		echo "Private key saved in $(PRIVATE_KEY_FILE)"; \
	else \
		echo "Private key already exists: $(PRIVATE_KEY_FILE)"; \
	fi

	@if [ ! -f $(PUBLIC_KEY_FILE) ]; then \
		echo "Extracting public key from the private key..."; \
		openssl ec -in $(PRIVATE_KEY_FILE) -pubout -out $(PUBLIC_KEY_FILE); \
		echo "Public key saved in $(PUBLIC_KEY_FILE)"; \
	else \
		echo "Public key already exists: $(PUBLIC_KEY_FILE)"; \
	fi

lint:
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Linter passed successfully"

start:
	@echo "Running API..."
	go run cmd/api/main.go

	
migrations:
	@echo "Runnig migrations..."
	go run database/migrations/main.go
	@echo "migrations executed succefully"

e2e:
	@echo "Runnig e2e tests"

.PHONY: all generate-keys lint clean