default:
    @just --list

build:
    go generate ./...
    GOOS=linux GOARCH=amd64 go build -o bin/lester main.go

generate-mocks:
    mockery --all --inpackage

test-coverage:
    go test -coverprofile=coverage.out ./... ;  go tool cover -html=coverage.out

test-watch:
    ls **/*.go | entr -c go test ./... -v
