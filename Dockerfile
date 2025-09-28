# Multi-stage build para otimizar o tamanho da imagem final
FROM golang:1.23.5-alpine AS builder

# Instalar dependências necessárias para compilação
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências primeiro (para cache de layers)
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Imagem final minimalista
FROM alpine:latest

# Instalar ca-certificates para HTTPS e timezone data
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root para segurança
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /app

# Copiar o binário compilado da etapa anterior
COPY --from=builder /app/main .

# Criar diretórios necessários e copiar arquivos de migração
RUN mkdir -p ./internal/infra/outgoing/postgres/migrations
COPY --from=builder /app/internal/infra/outgoing/postgres/migrations ./internal/infra/outgoing/postgres/migrations

# Mudar propriedade dos arquivos para o usuário não-root
RUN chown -R appuser:appgroup /app

# Mudar para usuário não-root
USER appuser

# Expor porta da aplicação
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]
