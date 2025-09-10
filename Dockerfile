
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cashflow ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cashflow .

EXPOSE 8080

CMD ["./cashflow"]
