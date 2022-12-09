FROM golang:1.19.3-alpine as builder

RUN apk update && apk add git

RUN go install golang.org/x/tools/cmd/goimports@latest

COPY . /app

WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

# FROM alpine:latest

# RUN mkdir /app

# COPY --from=builder /app/brokerApp /app

# CMD ["/app/brokerApp"]
