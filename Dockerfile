# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copiar código-fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bin/book-suggest-ai cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Instalar runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Copiar executável do builder
COPY --from=builder /app/bin/book-suggest-ai .

# Copiar templates e static files
COPY --from=builder /app/web ./web

# Copiar .env (opcional, pode ser overridden)
COPY .env .

# Criar diretório de dados
RUN mkdir -p data

# Expor porta
EXPOSE 8080

# Executar aplicação
CMD ["./book-suggest-ai"]
