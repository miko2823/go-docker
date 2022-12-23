FROM golang:1.19.3-alpine as builder

RUN apk update && apk add git make protobuf-dev

RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY . /app

WORKDIR /app

COPY go.* ./

RUN go mod download

RUN go mod verify

RUN go install -v github.com/rogpeppe/godef@v1.1.2

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

# FROM alpine:latest

# RUN mkdir /app

# COPY --from=builder /app/brokerApp /app

# CMD ["/app/brokerApp"]
