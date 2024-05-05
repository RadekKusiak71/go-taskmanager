build:
	@go build -o bin/go-taskmanager

run: build
	@./bin/go-taskmanager