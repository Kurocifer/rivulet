build:
	@cd rivulet-base && go build -o bin/rvt

run-base: build
	@./bin/rvt

test-base:
	@cd rivulet-base && go test ./... -v