FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o cashflow ./cmd/api

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/cashflow .
COPY .env .
EXPOSE 8080
CMD ["./cashflow"]
