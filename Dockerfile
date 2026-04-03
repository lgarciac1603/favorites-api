FROM golang:1.21-alpine AS builder

# Instalar git
RUN apk add --no-cache git

WORKDIR /app

# Copiar modulos
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Build optimizado
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o favorites-api .

# Runtime stage
FROM alpine:3.21

# Instalar certificados para HTTPS
RUN apk add --no-cache ca-certificates

# Crear usuario no-root
RUN adduser -D -h /home/appuser appuser

WORKDIR /app

# Copiar binario
COPY --from=builder /app/favorites-api /app/favorites-api

# Cambiar propietario
RUN chown -R appuser:appuser /app

EXPOSE 8090

USER appuser

CMD ["/app/favorites-api"]