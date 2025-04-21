FROM go:alpine1-22-alpine AS builder

RUN mkdir /app


COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o logger-service ./cmd/api 

# stage 2

FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/logger-service /app

CMD["/app/logger-service"]