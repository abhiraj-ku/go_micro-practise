# base go image :: Multi stage docker build
FROM golang:1.24-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o broker-service ./cmd/api

## stage 2

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/broker-service /app

CMD ["/app/broker-service"]