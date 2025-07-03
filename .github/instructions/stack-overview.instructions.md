---
description: Guia do stack principal de dependências Go e melhores práticas para o projeto mystery-gifter-api
applyTo: "**"
---
# Stack de Dependências Go — mystery-gifter-api

Este projeto utiliza um stack moderno para APIs REST em Go, com foco em produtividade, segurança e testabilidade. Abaixo estão as principais bibliotecas, suas versões e recomendações de uso.

## Principais Dependências

- **Go**: `1.23.5`
  - Use sempre a versão especificada no `go.mod` para evitar incompatibilidades.

### Web Framework & Middleware

- **github.com/gofiber/fiber/v2**: `v2.52.6`
  - Framework web rápido e minimalista.
  - Use middlewares oficiais para autenticação, CORS e logging.
- **github.com/gofiber/contrib/jwt**: `v1.0.10`
  - Middleware JWT para autenticação.
  - Prefira sempre validar algoritmos e claims explicitamente.

### Validação e Internacionalização

- **github.com/go-playground/validator/v10**: `v10.23.0`
  - Validação robusta de structs e campos.
  - Combine com `universal-translator` para mensagens customizadas.
- **github.com/go-playground/universal-translator**: `v0.18.1`
- **github.com/go-playground/locales**: `v0.14.1`
  - Suporte a internacionalização de mensagens de erro.

### Autenticação e Segurança

- **github.com/golang-jwt/jwt/v5**: `v5.2.1`
  - Use sempre a versão 5+ para evitar vulnerabilidades antigas.
  - Prefira algoritmos seguros (ex: HS256, RS256).
- **golang.org/x/crypto**: `v0.29.0`
  - Utilize funções modernas de hash e criptografia (ex: bcrypt, scrypt).

### Banco de Dados & Migrations

- **github.com/jmoiron/sqlx**: `v1.4.0`
  - Extensão para `database/sql` com suporte a named queries.
  - Use contextos (`context.Context`) em todas as queries.
- **github.com/lib/pq**: `v1.10.9`
  - Driver PostgreSQL.
  - Sempre feche conexões e use pooling.
- **github.com/golang-migrate/migrate/v4**: `v4.18.1`
  - Ferramenta de migração de banco de dados.
  - Mantenha scripts de migração versionados em `internal/outgoing/postgres/migrations/`.

### Utilitários

- **github.com/google/uuid**: `v1.6.0`
  - Para geração de UUIDs seguros.
- **github.com/joho/godotenv**: `v1.5.1`
  - Carregamento de variáveis de ambiente em desenvolvimento.

### Testes e Mocks

- **github.com/stretchr/testify**: `v1.9.0`
  - Testes unitários e assertions.
- **go.uber.org/mock**: `v0.5.0`
  - Geração de mocks para interfaces.

## Boas Práticas Gerais

- Sempre mantenha as versões travadas no `go.mod` para builds reprodutíveis.
- Atualize dependências críticas de segurança regularmente.
- Use contextos (`context.Context`) em handlers, serviços e queries.
- Separe camadas de domínio, aplicação e infraestrutura conforme a estrutura do projeto.
- Utilize builders e mocks para facilitar testes, conforme padrão do repositório.

## Exemplo de Importação Correta

```go
import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	// ... outros imports
)
```

> Consulte este arquivo para garantir que novas dependências estejam alinhadas com o stack e as melhores práticas do projeto.
