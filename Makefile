.PHONY: all deps fmt vet test builddev develop

EXECUTABLE ?= cron-service
IMAGE ?= fixate/$(EXECUTABLE):latest
COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: deps build test

deps:
	dep ensure

fmt:
	go fmt $(PACKAGES)

vet:
	go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

$(EXECUTABLE): $(wildcard *.go)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)'

docker: $(EXECUTABLE)
	docker build -t $(IMAGE) .

builddev:
	go build -ldflags '-s -w $(LDFLAGS)'

develop:
	realize start --build

runexample:
	./cron-service -m examples/manifest.yml --project-id digibet-dev-186112 --ensure-topics-created

build: $(EXECUTABLE)
