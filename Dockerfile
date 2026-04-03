FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o favorites-api .

FROM alpine:3.21

RUN adduser -D appuser
WORKDIR /app

COPY --from=builder /app/favorites-api /app/favorites-api

EXPOSE 8090

USER appuser
CMD ["/app/favorites-api"]
