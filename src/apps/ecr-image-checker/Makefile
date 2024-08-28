# Makefile for releasing ecr-image-checker
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=ecr-image-checker
DOCKER_REPOSITORY:=000000000000.dkr.ecr.us-east-1.amazonaws.com
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
IMG_PLATFORM ?= linux/$(shell go env GOARCH)
GIT_COMMIT:=$(shell git describe --dirty --always)
EXTRA_RUN_ARGS?=

run:
	go run -ldflags "-s -w -X github.com/rodrigorfk/k8s-kyverno-mutating-policies/ecr-image-checker/pkg/version.REVISION=$(GIT_COMMIT)" cmd/* \
	--level=debug $(EXTRA_RUN_ARGS)

.PHONY: test
test:
	go test ./... -coverprofile cover.out

build:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/rodrigorfk/k8s-kyverno-mutating-policies/ecr-image-checker/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/ecr-image-checker ./cmd/*

tidy:
	rm -f go.sum; go mod tidy -compat=1.22

vet:
	go vet ./...

fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

build-container:
	KO_DOCKER_REPO=$(DOCKER_IMAGE_NAME) \
	GIT_COMMIT=$(GIT_COMMIT) \
	ko build --tags $(TAG) --platform=${IMG_PLATFORM} cmd/ecr-image-checker.go -L --bare

