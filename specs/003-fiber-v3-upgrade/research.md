# Research: Atualização Go 1.26.4 + Fiber V2 → V3

**Feature**: 003-fiber-v3-upgrade  
**Data**: 2026-06-21  

---

## Decisão 1: Atualização do toolchain Go

**Decision**: Atualizar a diretiva `go` em `go.mod` de `1.25.1` para `1.26.4`.

**Rationale**: Go 1.26.4 é retrocompatível. A diretiva `go` em `go.mod` declara a versão mínima exigida — atualizar não exige mudanças no código-fonte. Go 1.26.x satisfaz o requisito mínimo do Fiber V3 (Go 1.25+).

**Alternatives considered**: Manter Go 1.25.1 — rejeitado pois o objetivo explícito da US1 é a atualização do toolchain.

**Ação**: `go 1.25.1` → `go 1.26.4` em `go.mod`, seguido de `go mod tidy`.

---

## Decisão 2: Paths de import do Fiber V3 e middlewares

**Decision**:

| Pacote (V2) | Pacote (V3) |
|-------------|-------------|
| `github.com/gofiber/fiber/v2` | `github.com/gofiber/fiber/v3` |
| `github.com/gofiber/fiber/v2/middleware/cors` | `github.com/gofiber/fiber/v3/middleware/cors` |
| `github.com/gofiber/fiber/v2/middleware/recover` | `github.com/gofiber/fiber/v3/middleware/recover` |
| `github.com/gofiber/contrib/jwt` | `github.com/gofiber/contrib/v3/jwt` |

**Rationale**: Os caminhos de middleware seguem o mesmo padrão de prefixo que em V2; apenas o módulo raiz muda de `v2` para `v3`. O pacote JWT do contrib tem módulo próprio com versão `v3`.

**Alternatives considered**: N/A — caminhos definidos pelo upstream do Fiber.

---

## Decisão 3: Assinatura de handler

**Decision**: Handlers mudam de `func(ctx *fiber.Ctx) error` para `func(ctx fiber.Ctx) error`.

`fiber.Ctx` em V3 é uma interface (não mais um ponteiro para struct). Todos os controllers, error handlers e middleware ErrorHandlers precisam ser atualizados.

**Rationale**: Mudança imposta pelo Fiber V3 — necessária para compilar.

**Impacto em arquivos de produção**:
- `internal/infra/entrypoint/middlewares.go` — `ErrorHandler`
- `internal/infra/entrypoint/error_handler.go` — `CustomErrorHandler` e `sendError`
- `internal/infra/entrypoint/rest/*.go` — todos os métodos de controller

**Impacto em arquivos de teste**:
- Handler lambdas inline nos `*_test.go` dos pacotes `entrypoint` e `rest`

---

## Decisão 4: API de binding de request

**Decision**:

| V2 | V3 |
|----|-----|
| `ctx.BodyParser(&dto)` | `ctx.Bind().Body(&dto)` |
| `ctx.QueryParser(&dto)` | `ctx.Bind().Query(&dto)` |

Ambos retornam `error` ao falhar — o padrão de tratamento permanece idêntico:

```go
if err := ctx.Bind().Body(&dto); err != nil {
    return fiber.NewError(fiber.StatusUnprocessableEntity)
}
```

**Rationale**: API de binding unificada do Fiber V3; comportamento de erro equivalente para JSON malformado.

**Alternatives considered**: Manter `BodyParser` via camada de compatibilidade — rejeitado pois a função foi removida no V3.

---

## Decisão 5: Acesso ao token JWT autenticado

**Decision**: Substituir `ctx.Locals("user")` por `jwtware.FromContext(ctx)` nos controllers.

`jwtware.FromContext(ctx fiber.Ctx) *jwt.Token` é a função tipada exposta pelo `gofiber/contrib/v3/jwt`. Elimina a dependência da string `"user"` como chave e retorna o token com tipo conhecido.

A interface `AuthTokenManager.GetAuthUserID(token any)` e sua implementação em `JWTAuthTokenManager` **não mudam** — `*jwt.Token` satisfaz `any`.

**Rationale**: `jwtware.FromContext` é a API canônica do V3 para recuperar o token após validação pelo middleware. Mais segura que acesso por string.

**Alternatives considered**: `fiber.Locals[*jwt.Token](ctx, "user")` — possível, mas depende do nome interno da chave usada pelo middleware, que pode mudar. `FromContext` é a API pública estável.

---

## Decisão 6: Configuração CORS

**Decision**: O `cors.New()` sem argumentos no `runner.go` (config padrão) continua sendo chamado da mesma forma. Nenhuma mudança funcional.

Se no futuro a config for customizada, `AllowOrigins`, `AllowMethods` e `AllowHeaders` devem usar `[]string` (não strings separadas por vírgula como no V2).

**Rationale**: Comportamento padrão (wildcard) é mantido pelo middleware V3.

---

## Decisão 7: Middleware recover

**Decision**: `recover.New()` sem argumentos continua sendo chamado da mesma forma. Nenhuma mudança funcional.

**Rationale**: API padrão preservada no V3.

---

## Decisão 8: go-swagger e geração de documentação

**Decision**: As anotações Swagger (`// @Summary`, `// @Param`, etc.) são comentários independentes do framework. O `go-swagger` as escaneia no código-fonte sem conhecimento do Fiber. Nenhuma mudança nas anotações é necessária.

**Rationale**: A ferramenta `go-swagger` trabalha com AST do Go, não com tipos do Fiber. A mudança de `*fiber.Ctx` para `fiber.Ctx` não afeta o parsing das anotações.

**Ação**: Executar `make generate-docs` ao final da migração para confirmar que a geração continua funcionando.

---

## Decisão 9: Escopo de atualização de outras dependências

**Decision**: Atualizar somente Go + Fiber + JWT contrib como alvo primário. Outras dependências (`sqlx`, `squirrel`, `golang-jwt`, `caarlos0/env`, etc.) serão atualizadas apenas se falharem a compilação com as novas versões.

**Rationale**: Minimiza superfície de risco e facilita isolamento de regressões.

---

## Decisão 10: Atualização da Constituição do projeto

**Decision**: A Constituição contém referências a padrões Fiber V2 que precisam ser atualizadas:

| Seção | Conteúdo atual (V2) | Novo conteúdo (V3) |
|-------|---------------------|---------------------|
| Princípio IV — assinatura de handler | `func (c *XController) Method(ctx *fiber.Ctx) error` | `func (c *XController) Method(ctx fiber.Ctx) error` |
| Princípio IV — fluxo de requisição | `BodyParser →` | `Bind().Body() →` |
| Princípio IV — falha de parsing | `Falha em BodyParser DEVE retornar...` | `Falha em Bind().Body() DEVE retornar...` |
| Princípio IV — extração de userID | `ctx.Locals("user")` | `jwtware.FromContext(ctx)` |
| Padrões Tecnológicos | `Fiber v2` | `Fiber v3` |

**Rationale**: A Constituição é a fonte canônica de padrões — deve refletir os padrões V3 após a migração para guiar implementações futuras corretamente.
