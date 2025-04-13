# base go image :: Multi stage docker build
FROM golang:1.24-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o auth-service ./cmd/api

## stage 2

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/auth-service /app

CMD ["/app/auth-service"]