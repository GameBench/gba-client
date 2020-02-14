build:
	GOOS=darwin GOARCH=amd64 go build -trimpath -o gba-client-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -trimpath -o gba-client-linux-amd64
	GOOS=windows GOARCH=amd64 go build -trimpath -o gba-client-windows-amd64

upload:
	gsutil cp gba-client-darwin-amd64 gs://release-application-artifacts/gba-client/
	gsutil cp gba-client-linux-amd64 gs://release-application-artifacts/gba-client/
	gsutil cp gba-client-windows-amd64 gs://release-application-artifacts/gba-client/

.PHONY: build upload
