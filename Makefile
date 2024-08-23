.PHONY: build dev run

build:
	@echo "Building..."
	docker build -t naabu-http-wrapper .

dev:
	@echo "Starting dev server..."
	@go run .

run:
	@echo "Starting docker container..."
	@docker run -p 8080:8080 naabu-http-wrapper 