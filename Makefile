# rivulet base build section.
buildbase:
	@cd rivulet-base && go build -o ../bin/base/rvt

runbase: buildbase
	@./bin/base/rvt

testbase:
	@cd rivulet-base && go test ./... -v


# rivulet cli build section.
buildcli:
	@cd rivulet-cli && go build -o ../bin/cli/rvt-cli