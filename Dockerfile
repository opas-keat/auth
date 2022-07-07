# ./Dockerfile
FROM golang:1.18.3-alpine3.15 AS builder
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container.
ADD cmd/ cmd/
ADD .env .env

# Set necessary environment variables needed for our image 
# and build the API server.
RUN go build -o main cmd/v1/main.go

FROM alpine:latest
# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist/
# Copy binary from build to dist folder
COPY --from=builder /build .
# Export necessary port
EXPOSE 8080
# Command to run when starting the container.
ENTRYPOINT ["./main"]
