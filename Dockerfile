FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/migrations /migrations

EXPOSE ${APP_PORT}

CMD ["./main"]
