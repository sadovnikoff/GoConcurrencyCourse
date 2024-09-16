
CLI_APP_NAME=database-cli

build-cli:
	go build -o ${CLI_APP_NAME} cmd/cli/main.go

run-cli: build-cli
	./${CLI_APP_NAME} $(ARGS)

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out