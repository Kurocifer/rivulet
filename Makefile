build:
	@cd rivulet-base && go build -o ../bin/base/rvt

runbase: build
	@./bin/base/rvt

testbase:
	@cd rivulet-base && go test ./... -v