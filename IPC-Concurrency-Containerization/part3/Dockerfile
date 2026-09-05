FROM golang:1.22-alpine AS builder

ENV GOPROXY=https://mirror-go.runflare.com

WORKDIR /app

COPY . .

RUN go build -o service main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8080

CMD ["./service"]
