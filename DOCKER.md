# Mystery Gifter API - Docker Setup

Este documento explica como executar a aplicação Mystery Gifter API usando Docker.

## Pré-requisitos

- Docker
- Docker Compose

## Configuração

1. **Copie o arquivo de exemplo de variáveis de ambiente:**
   ```bash
   cp .env.example .env
   ```

2. **Edite o arquivo `.env` com suas configurações:**
   ```bash
   # Configurações do Banco de Dados
   DB_HOST=localhost
   DB_PORT=5432
   DB_DATABASE=mystery_gifter
   DB_USERNAME=postgres
   DB_PASSWORD=postgres

   # Configurações de Autenticação
   AUTH_SECRET_KEY=your-secret-key-here-change-in-production
   AUTH_SESSION_DURATION=24h
   ```

## Executando a Aplicação

### Desenvolvimento (com hot reload)

```bash
# Construir e executar todos os serviços
docker-compose up --build

# Executar em background
docker-compose up -d --build
```

### Produção

```bash
# Construir e executar em modo produção
docker-compose -f docker-compose.yml up --build -d
```

## Comandos Úteis

### Ver logs
```bash
# Todos os serviços
docker-compose logs

# Apenas a API
docker-compose logs api

# Apenas o banco de dados
docker-compose logs db

# Seguir logs em tempo real
docker-compose logs -f api
```

### Parar os serviços
```bash
# Parar todos os serviços
docker-compose down

# Parar e remover volumes (CUIDADO: apaga dados do banco)
docker-compose down -v
```

### Reconstruir apenas a API
```bash
docker-compose build api
docker-compose up -d api
```

### Acessar o container da API
```bash
docker-compose exec api sh
```

### Acessar o banco de dados
```bash
docker-compose exec db psql -U postgres -d mystery_gifter
```

## Portas

- **API**: http://localhost:8080
- **Banco de Dados**: localhost:5432

## Health Checks

A aplicação inclui health check apenas para o banco de dados:

- **Banco de Dados**: Verifica se o PostgreSQL está pronto para conexões

## Estrutura dos Serviços

### API Service
- **Imagem**: Construída localmente usando o Dockerfile
- **Porta**: 8080
- **Dependências**: Aguarda o banco de dados estar saudável
- **Variáveis de Ambiente**: Configuradas via arquivo `.env`

### Database Service
- **Imagem**: postgres:latest
- **Porta**: 5432
- **Volume**: Persistência de dados em `db_data`
- **Health Check**: Verifica conectividade com PostgreSQL

## Troubleshooting

### Problema: API não consegue conectar ao banco
- Verifique se o `DB_HOST` está configurado como `db` (nome do serviço)
- Verifique se o banco está saudável: `docker-compose ps`

### Problema: Porta já está em uso
- Pare outros serviços que possam estar usando as portas 8080 ou 5432
- Ou altere as portas no `docker-compose.yml`

### Problema: Erro de permissão
- Certifique-se de que o Docker tem permissões adequadas
- No Linux: `sudo usermod -aG docker $USER` e faça logout/login

## Desenvolvimento

Para desenvolvimento local sem Docker:

1. Execute apenas o banco: `docker-compose up db -d`
2. Configure `DB_HOST=localhost` no `.env`
3. Execute a aplicação localmente: `go run cmd/api/main.go`
