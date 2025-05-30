FROM golang:1.24-alpine AS builder

RUN mkdir /app

# should copy the go mod file and install 
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o listener-service ./cmd/api


# stage 2

FROM alpine:latest
WORKDIR /app


COPY --from=builder /app/listener-service /app

CMD["/app/listener-service"]