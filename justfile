default:
    @just --list

# build mqtt-blackbox-exporter binary
build:
    go build -o mqtt-blackbox-exporter ./cmd/mqtt-blackbox-exporter

# update go packages
update:
    @cd ./cmd/mqtt-blackbox-exporter && go get -u

# run tests
test:
    go test -v ./... -covermode=atomic -coverprofile=coverage.out

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml
