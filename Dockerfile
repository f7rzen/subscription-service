FROM golang:1.26.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/subscription-service ./cmd/app


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/subscription-service .

EXPOSE 8080

CMD ["./subscription-service"]