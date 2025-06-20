build:
	@go build -o bin/rvt

run: build
	@./bin/rvt

test:
	@go test ./... -v