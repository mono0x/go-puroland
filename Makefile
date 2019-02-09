GO=go
GOBIN=$(PWD)/bin
TESTOPTS=-coverprofile=result.coverprofile -v -race ./...

all: deps test

setup:
	GOBIN=$(GOBIN) GO111MODULE=on go install honnef.co/go/tools/cmd/staticcheck

deps:
	GO111MODULE=on go mod tidy

test:
	GO111MODULE=on $(GO) mod verify
	GO111MODULE=on $(GO) vet ./...
	GO111MODULE=on $(GO) test $(TESTOPTS)
	GO111MODULE=on PATH=$(GOBIN):$(PATH) staticcheck ./...
