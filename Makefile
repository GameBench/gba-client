GIT_HASH := $(shell git rev-parse HEAD)
VERSION := 0.1.0

build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.commitHash=$(GIT_HASH) -X main.version=$(VERSION) -s" -trimpath -o gba-client-darwin-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.commitHash=$(GIT_HASH) -X main.version=$(VERSION)" -trimpath -o gba-client-linux-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.commitHash=$(GIT_HASH) -X main.version=$(VERSION)" -trimpath -o gba-client-windows-amd64

upload:
	gsutil -h x-goog-meta-commit-hash:$(GIT_HASH) cp gba-client-darwin-amd64 gs://release-application-artifacts/gba-client/
	gsutil -h x-goog-meta-commit-hash:$(GIT_HASH) cp gba-client-linux-amd64 gs://release-application-artifacts/gba-client/
	gsutil -h x-goog-meta-commit-hash:$(GIT_HASH) cp gba-client-windows-amd64 gs://release-application-artifacts/gba-client/

.PHONY: build upload
