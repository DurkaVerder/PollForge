FROM golang:1.23.8 AS builder

WORKDIR /app

COPY . .


RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/app/main.go




FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

RUN ls -l

CMD ["./app"]