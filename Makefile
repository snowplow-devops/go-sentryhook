.PHONY: all format lint tidy

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all:
	GO111MODULE=on go build .

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	GO111MODULE=on go fmt .
	GO111MODULE=on gofmt -s -w .

lint:
	GO111MODULE=on go get -u golang.org/x/lint/golint
	GO111MODULE=on golint .

tidy:
	GO111MODULE=on go mod tidy
