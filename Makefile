hello:
	echo "Hello"

build:
	GOOS=windows GOARCH=amd64 go build -o bin/auth/authapi.exe cmd/v1/main.go

run:
	go run cmd/main.go
