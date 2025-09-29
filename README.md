# 🎁 Mystery Gifter API

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52.6-00ADD8.svg)](https://gofiber.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Latest-336791.svg)](https://postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED.svg)](https://docker.com/)
[![Swagger](https://img.shields.io/badge/Swagger-Documented-85EA2D.svg)](http://localhost:8080/swagger/)

Uma API REST moderna desenvolvida em Go para gerenciar grupos de **Secret Santa** (Amigo Secreto). Permite criar grupos, adicionar participantes e gerar aleatoriamente os pares de quem presenteia quem.

## 📋 Índice

- [Características](#-características)
- [Tecnologias](#-tecnologias)
- [Arquitetura](#-arquitetura)
- [Pré-requisitos](#-pré-requisitos)
- [Instalação e Execução](#-instalação-e-execução)
- [Configuração](#-configuração)
- [Documentação da API](#-documentação-da-api)
- [Endpoints](#-endpoints)
- [Exemplos de Uso](#-exemplos-de-uso)
- [Desenvolvimento](#-desenvolvimento)
- [Testes](#-testes)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Contribuição](#-contribuição)
- [Licença](#-licença)

## ✨ Características

- 🎯 **Gestão Completa de Grupos**: Criar, gerenciar e arquivar grupos de Secret Santa
- 👥 **Sistema de Usuários**: Cadastro, autenticação e busca de usuários
- 🎲 **Geração Aleatória de Matches**: Algoritmo seguro para criar pares de presenteio
- 🔐 **Autenticação JWT**: Sistema seguro de autenticação com tokens
- 📊 **Busca e Filtros**: Sistema avançado de busca com paginação
- 🏗️ **Arquitetura Limpa**: Implementação seguindo Clean Architecture
- 🧪 **Cobertura de Testes**: Testes unitários abrangentes
- 📚 **Documentação Swagger**: API completamente documentada
- 🐳 **Docker Ready**: Containerização completa para desenvolvimento e produção

## 🛠️ Tecnologias

### Backend
- **[Go 1.25.1](https://golang.org/)** - Linguagem principal
- **[Fiber v2.52.6](https://gofiber.io/)** - Framework web de alta performance
- **[PostgreSQL](https://postgresql.org/)** - Banco de dados relacional
- **[JWT](https://jwt.io/)** - Autenticação baseada em tokens
- **[BCrypt](https://en.wikipedia.org/wiki/Bcrypt)** - Hash de senhas seguro

### Ferramentas de Desenvolvimento
- **[Docker](https://docker.com/)** - Containerização
- **[Swagger/OpenAPI](https://swagger.io/)** - Documentação da API
- **[Testify](https://github.com/stretchr/testify)** - Framework de testes
- **[Mock](https://github.com/uber-go/mock)** - Geração de mocks
- **[Migrate](https://github.com/golang-migrate/migrate)** - Migrações de banco

## 🏗️ Arquitetura

O projeto segue os princípios da **Clean Architecture** (Arquitetura Hexagonal), garantindo:

- **Separação de responsabilidades** clara
- **Testabilidade** alta
- **Flexibilidade** para mudanças
- **Manutenibilidade** otimizada

```
📁 internal/
├── 🎯 domain/          # Entidades e regras de negócio
├── 🔧 application/      # Casos de uso e serviços
└── 🏗️ infra/           # Infraestrutura e adaptadores
    ├── config/         # Configurações
    ├── entrypoint/     # Controllers REST e rotas
    └── outgoing/        # Repositórios e serviços externos
```

## 📋 Pré-requisitos

- **[Docker](https://docs.docker.com/get-docker/)** (versão 20.10+)
- **[Docker Compose](https://docs.docker.com/compose/install/)** (versão 2.0+)

> 💡 **Dica**: Para desenvolvimento local, você também pode usar Go 1.25.1+ instalado diretamente no sistema.

## 🚀 Instalação e Execução

### 1. Clone o repositório

```bash
git clone https://github.com/waliqueiroz/mystery-gifter-api.git
cd mystery-gifter-api
```

### 2. Configure as variáveis de ambiente

Copie o arquivo de exemplo e configure suas variáveis:

```bash
cp .env.example .env
```

Edite o arquivo `.env` com suas configurações:

```bash
# Configurações do Banco de Dados
DB_HOST=db
DB_PORT=5432
DB_DATABASE=mystery_gifter
DB_USERNAME=postgres
DB_PASSWORD=postgres

# Configurações de Autenticação
AUTH_SECRET_KEY=your-super-secret-key-change-in-production-minimum-32-chars
AUTH_SESSION_DURATION=24h
```

> ⚠️ **Importante**: Altere a `AUTH_SECRET_KEY` para uma chave segura em produção!

### 3. Execute com Docker Compose

```bash
# Desenvolvimento (com rebuild automático)
docker-compose up --build

# Executar em background
docker-compose up -d --build

# Apenas produção
docker-compose up -d
```

### 4. Verifique se está funcionando

- **API**: http://localhost:8080
- **Banco de Dados**: localhost:5432

### 5. Gerar e visualizar a documentação Swagger

A documentação Swagger precisa ser gerada e servida separadamente:

```bash
# Gerar a documentação Swagger
make generate-docs

# Servir a documentação (em uma nova aba do terminal)
make serve-docs
```

Após executar `make serve-docs`, acesse:
- **Documentação Swagger**: http://localhost:8081

## ⚙️ Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão | Obrigatória |
|----------|-----------|--------|-------------|
| `DB_HOST` | Host do banco de dados | `db` | ✅ |
| `DB_PORT` | Porta do banco de dados | `5432` | ✅ |
| `DB_DATABASE` | Nome do banco | `mystery_gifter_db` | ✅ |
| `DB_USERNAME` | Usuário do banco | `postgres` | ✅ |
| `DB_PASSWORD` | Senha do banco | - | ✅ |
| `AUTH_SECRET_KEY` | Chave secreta para JWT | - | ✅ |
| `AUTH_SESSION_DURATION` | Duração da sessão | `24h` (apenas no Docker) | ✅ |

> ⚠️ **Nota**: `AUTH_SESSION_DURATION` é obrigatória. No Docker Compose há um valor padrão (`24h`), mas para execução local você deve defini-la explicitamente.

### Estados dos Grupos

- **`OPEN`**: Aceitando novos usuários
- **`MATCHED`**: Matches já gerados
- **`ARCHIVED`**: Grupo arquivado (não pode ser reaberto)

## 📚 Documentação da API

A API está completamente documentada com **Swagger/OpenAPI**. Para visualizar a documentação:

### Gerar e Servir a Documentação

```bash
# 1. Gerar a documentação a partir dos comentários no código
make generate-docs

# 2. Servir a documentação interativa
make serve-docs
```

### Acessar a Documentação

Após executar `make serve-docs`, acesse:
- **Swagger UI**: http://localhost:8081
- **Especificação OpenAPI**: http://localhost:8081/swagger.json

> 💡 **Dica**: A documentação é gerada automaticamente a partir dos comentários `swagger:` no código. Sempre atualize os comentários quando modificar a API.

## 🔗 Endpoints

### 🔐 Autenticação
- `POST /api/v1/login` - Login e obtenção de token JWT

### 👥 Usuários
- `POST /api/v1/users` - Criar novo usuário
- `GET /api/v1/users` - Buscar usuários (com filtros e paginação)
- `GET /api/v1/users/{id}` - Obter usuário por ID

### 🎁 Grupos
- `GET /api/v1/groups` - Buscar grupos (com filtros e paginação)
- `POST /api/v1/groups` - Criar novo grupo
- `GET /api/v1/groups/{id}` - Obter grupo por ID
- `POST /api/v1/groups/{id}/users` - Adicionar usuário ao grupo
- `DELETE /api/v1/groups/{id}/users/{userId}` - Remover usuário do grupo
- `POST /api/v1/groups/{id}/matches` - Gerar matches aleatórios
- `GET /api/v1/groups/{id}/matches/user` - Obter match do usuário logado
- `POST /api/v1/groups/{id}/reopen` - Reabrir grupo
- `POST /api/v1/groups/{id}/archive` - Arquivar grupo

> 🔒 **Nota**: Todos os endpoints exceto `POST /api/v1/users` e `POST /api/v1/login` requerem autenticação JWT.

## 💡 Exemplos de Uso

### 1. Criar um usuário

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João",
    "surname": "Silva",
    "email": "joao@example.com",
    "password": "minhasenha123",
    "password_confirm": "minhasenha123"
  }'
```

### 2. Fazer login

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "minhasenha123"
  }'
```

### 3. Criar um grupo (com token JWT)

```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_JWT_TOKEN_AQUI" \
  -d '{
    "name": "Secret Santa 2024"
  }'
```

### 4. Gerar matches

```bash
curl -X POST http://localhost:8080/api/v1/groups/{group_id}/matches \
  -H "Authorization: Bearer SEU_JWT_TOKEN_AQUI"
```

## 🛠️ Desenvolvimento

### Executar localmente (sem Docker)

1. **Inicie apenas o banco de dados:**
   ```bash
   docker-compose up db -d
   ```

2. **Configure o `.env` para desenvolvimento local:**
   ```bash
   DB_HOST=localhost
   # ... outras configurações
   ```

3. **Execute a aplicação:**
   ```bash
   go run cmd/api/main.go
   ```

### Comandos úteis do Docker

```bash
# Ver logs em tempo real
docker-compose logs -f api

# Parar todos os serviços
docker-compose down

# Reconstruir apenas a API
docker-compose build api && docker-compose up -d api

# Acessar o container da API
docker-compose exec api sh

# Acessar o banco de dados
docker-compose exec db psql -U postgres -d mystery_gifter
```

### Comandos úteis do Makefile

```bash
# Ver todos os comandos disponíveis
make help

# Gerar documentação Swagger
make generate-docs

# Servir documentação Swagger
make serve-docs

# Executar testes
make test

# Compilar a aplicação
make build

# Executar a aplicação localmente
make run

# Instalar ferramentas necessárias
make install-tools
```

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Ou usando go diretamente
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes de um pacote específico
go test ./internal/domain/...

# Executar testes com verbose
go test -v ./...
```

## 📁 Estrutura do Projeto

```
mystery-gifter-api/
├── 📁 cmd/api/                    # Ponto de entrada da aplicação
├── 📁 internal/
│   ├── 📁 application/            # Camada de aplicação (casos de uso)
│   ├── 📁 domain/                 # Camada de domínio (entidades e regras)
│   └── 📁 infra/                  # Camada de infraestrutura
│       ├── 📁 config/             # Configurações
│       ├── 📁 entrypoint/         # Controllers REST e rotas
│       └── 📁 outgoing/          # Repositórios e serviços externos
├── 📁 pkg/                       # Pacotes reutilizáveis
├── 📁 test/                      # Utilitários de teste
├── 📁 docs/                      # Documentação Swagger
├── 🐳 docker-compose.yml         # Configuração Docker Compose
├── 🐳 Dockerfile                 # Imagem Docker da aplicação
├── 📄 go.mod                     # Dependências Go
└── 📄 README.md                  # Este arquivo
```

## 🤝 Contribuição

Contribuições são bem-vindas! Para contribuir:

1. **Fork** o projeto
2. **Crie** uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. **Commit** suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. **Push** para a branch (`git push origin feature/AmazingFeature`)
5. **Abra** um Pull Request

### Padrões de Código

- Siga as convenções do Go
- Escreva testes para novas funcionalidades
- Mantenha a cobertura de testes alta
- Documente mudanças na API no Swagger
- Use commits semânticos

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

<div align="center">

**Desenvolvido com ❤️ em Go**

[⭐ Dê uma estrela](https://github.com/waliqueiroz/mystery-gifter-api) se este projeto te ajudou!

</div>
