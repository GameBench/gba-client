GIT_HASH := $(shell git rev-parse HEAD)

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.commitHash=$(GIT_HASH)" -trimpath -o gba-client-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.commitHash=$(GIT_HASH)" -trimpath -o gba-client-linux-amd64
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.commitHash=$(GIT_HASH)" -trimpath -o gba-client-windows-amd64

upload:
	gsutil -h x-goog-meta-commit-hash:$(GIT_HASH) cp gba-client-darwin-amd64 gs://release-application-artifacts/gba-client/
	gsutil -h x-goog-meta-commit-hash:$(GIT_HASH) cp gba-client-linux-amd64 gs://release-application-artifacts/gba-client/
	gsutil -h x-goog-meta-commit-hash:$(GIT_HASH) cp gba-client-windows-amd64 gs://release-application-artifacts/gba-client/

.PHONY: build upload
