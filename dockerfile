FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go build -o skii-db .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/skii-db /app/skii-db

RUN mkdir -p /app/data

VOLUME ["/app/data"]

ENV DATA_PATH=/app/data/data.txt

ENTRYPOINT ["/app/skii-db"]