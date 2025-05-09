FROM golang:1.24-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mail-service ./cmd/app

## stage 2

FROM alpine:latest


WORKDIR /app

COPY --from=builder /app/mail-service /app
COPY templates /templates

CMD["/app/mail-service"]