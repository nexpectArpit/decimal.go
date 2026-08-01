# Dockerfile for decimal.go port build & verification
FROM golang:1.22-alpine AS builder

RUN apk add --no-linux make nodejs npm git

WORKDIR /app

COPY . .

RUN make build
RUN make test

CMD ["make", "test"]
