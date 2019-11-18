build:
	GOOS=darwin GOARCH=amd64 go build -o gba-client-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -o gba-client-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o gba-client-windows-amd64

.PHONY: build