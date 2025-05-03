FROM golang:1.23.3 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app ./cmd/app/main.go




FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

CMD ["./app"]