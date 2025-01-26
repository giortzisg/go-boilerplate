FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

ARG CONFIG_ENV

COPY --from=builder /app/config/${CONFIG_ENV}.yaml /app/config/config.yaml

EXPOSE 8080

ENTRYPOINT ["./server", "-config", "config/config.yaml"]
