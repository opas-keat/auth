hello:
	echo "Hello"

build:
	GOOS=windows GOARCH=amd64 go build -o bin/ectapi/ectapi.exe cmd/v1/main.go

run:
	go run cmd/v1/main.go
