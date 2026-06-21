# Plano de Implementação: Atualização Go 1.26.4 + Fiber V2 → V3

**Branch**: `003-fiber-v3-upgrade` | **Data**: 2026-06-21 | **Spec**: [spec.md](spec.md)

## Resumo

Migração em dois checkpoints sequenciais: primeiro atualização do toolchain Go para 1.26.4 (pré-requisito isolado), depois migração do Fiber V2 para V3 com todas as breaking changes mapeadas. A API pública permanece inalterada — a migração é transparente para o usuário final. A Constituição do projeto também será atualizada para refletir os novos padrões V3.

## Contexto Técnico

**Linguagem/Versão**: Go 1.25.1 → Go 1.26.4  
**Framework Web**: Fiber v2.52.6 → Fiber v3 (latest)  
**JWT Middleware**: `github.com/gofiber/contrib/jwt v1.0.10` → `github.com/gofiber/contrib/v3/jwt`  
**Storage**: PostgreSQL via sqlx + squirrel (sem mudança)  
**Testes**: testify/assert + gomock (sem mudança)  
**Plataforma**: Linux server, porta 8080 (sem mudança)  
**Performance**: meta p95 < 200ms para CRUD padrão (sem mudança)  
**Escopo de deps adicionais**: somente se falharem compilação com Go 1.26.4 ou Fiber V3  

## Verificação da Constituição

| Princípio | Status | Observação |
|-----------|--------|------------|
| I — Clean Architecture | ✅ Conforme | Migration não toca domain/application |
| II — Disciplina de Testes | ✅ Conforme | Test files atualizados; todos os testes devem passar |
| III — Validação Orientada ao Domínio | ✅ N/A | Sem mudanças em validação |
| IV — Contrato de API Consistente | ⚠️ Requer atualização | Constituição referencia `*fiber.Ctx`, `BodyParser`, `ctx.Locals("user")` — serão atualizados como task obrigatória |
| V — Abstração de Infraestrutura | ✅ Conforme | Interfaces de domínio preservadas |
| VI — Simplicidade & YAGNI | ✅ Conforme | Sem novas abstrações |
| VII — Performance & Observabilidade | ✅ Conforme | Nenhum middleware de observabilidade alterado |
| VIII — Idioma dos Artefatos | ✅ Conforme | |
| Padrões Tecnológicos | ⚠️ Requer atualização | "Fiber v2" → "Fiber v3" na Constituição |

**Resolução dos gates ⚠️**: As referências desatualizadas na Constituição serão corrigidas como última task do Checkpoint 2. Não é uma violação arquitetural — é uma consequência esperada da migration que a própria task resolve.

## Estrutura do Projeto

### Documentação (esta feature)

```text
specs/003-fiber-v3-upgrade/
├── plan.md          # Este arquivo
├── research.md      # Decisões de pesquisa (gerado)
├── tasks.md         # Gerado por /speckit.tasks
└── checklists/
    └── requirements.md
```

### Código-fonte impactado

```text
go.mod / go.sum                                              # versão Go + deps Fiber

internal/infra/
├── runner.go                                                # imports CORS/recover
└── entrypoint/
    ├── middlewares.go                                       # ErrorHandler signature
    ├── routes.go                                            # import fiber
    ├── error_handler.go                                     # CustomErrorHandler + sendError
    └── rest/
        ├── auth_controller.go                               # BodyParser → Bind().Body()
        ├── user_controller.go                               # BodyParser + QueryParser
        ├── group_controller.go                              # BodyParser + QueryParser + Locals
        ├── group_invite_controller.go                       # Locals
        ├── auth_controller_test.go                          # handler signatures
        ├── user_controller_test.go                          # handler signatures
        ├── group_controller_test.go                         # handler signatures
        └── group_invite_controller_test.go                  # handler signatures

.specify/memory/constitution.md                              # Padrões Fiber V2 → V3
```

## Checkpoint 1: Atualização do Go para 1.26.4

**Escopo**: Apenas `go.mod` + verificação de compilação. Nenhum arquivo de código-fonte alterado.

**Entregável**: PR independente, branch stacked em `main`.

### Tarefas

#### T1.1 — Atualizar diretiva `go` em `go.mod`

Alterar `go 1.25.1` para `go 1.26.4` em `go.mod`.

#### T1.2 — Executar `go mod tidy`

Garante que `go.sum` e dependências indiretas estejam consistentes com a nova versão do toolchain. Se qualquer dependência falhar a compilação, atualizá-la para a versão mínima compatível (ver Decisão 9 em research.md).

#### T1.3 — Verificar build e testes

```bash
make build
make test
```

Ambos devem passar sem erros. Se houver falhas, corrigir antes de avançar.

---

## Checkpoint 2: Migração Fiber V2 → V3

**Escopo**: Atualização de dependências, adaptação de todo código infra que usa o Fiber, atualização de testes, atualização da Constituição.

**Pré-requisito**: Checkpoint 1 mergeado.

**Entregável**: PR stacked no Checkpoint 1.

### Tarefas

#### T2.1 — Atualizar dependências Fiber em `go.mod`

```bash
go get github.com/gofiber/fiber/v3@latest
go get github.com/gofiber/contrib/v3/jwt@latest
go mod tidy
```

Remove `github.com/gofiber/fiber/v2` e `github.com/gofiber/contrib/jwt` como dependências diretas.

---

#### T2.2 — Atualizar `internal/infra/runner.go`

**Mudanças**:
- Imports: `fiber/v2` → `fiber/v3`; `fiber/v2/middleware/cors` → `fiber/v3/middleware/cors`; `fiber/v2/middleware/recover` → `fiber/v3/middleware/recover`
- `fiber.Config{ErrorHandler: ...}` — verificar se a assinatura do `ErrorHandler` mudou (passa a aceitar `fiber.Ctx` em vez de `*fiber.Ctx`); a referência `entrypoint.CustomErrorHandler` já estará correta após T2.4

Nenhuma mudança funcional: `cors.New()` e `recover.New()` permanecem sem argumentos.

---

#### T2.3 — Atualizar `internal/infra/entrypoint/middlewares.go`

**Mudanças**:
- Import: `github.com/gofiber/contrib/jwt` → `github.com/gofiber/contrib/v3/jwt`
- Import: `github.com/gofiber/fiber/v2` → `github.com/gofiber/fiber/v3`
- `ErrorHandler`: assinatura `func(c *fiber.Ctx, err error) error` → `func(c fiber.Ctx, err error) error`

---

#### T2.4 — Atualizar `internal/infra/entrypoint/error_handler.go`

**Mudanças**:
- Import: `fiber/v2` → `fiber/v3`
- `CustomErrorHandler(ctx *fiber.Ctx, err error) error` → `CustomErrorHandler(ctx fiber.Ctx, err error) error`
- `sendError(ctx *fiber.Ctx, ...) error` → `sendError(ctx fiber.Ctx, ...) error`

---

#### T2.5 — Atualizar `internal/infra/entrypoint/routes.go`

**Mudanças**:
- Import: `fiber/v2` → `fiber/v3`
- Assinaturas dos parâmetros `fiber.Router` e `fiber.Handler` — permanecem com mesmo nome, apenas o pacote muda

---

#### T2.6 — Atualizar `internal/infra/entrypoint/rest/auth_controller.go`

**Mudanças**:
- Import: `fiber/v2` → `fiber/v3`
- `func (c *AuthController) Login(ctx *fiber.Ctx) error` → `func (c *AuthController) Login(ctx fiber.Ctx) error`
- `ctx.BodyParser(&credentialsDTO)` → `ctx.Bind().Body(&credentialsDTO)`

---

#### T2.7 — Atualizar `internal/infra/entrypoint/rest/user_controller.go`

**Mudanças**:
- Import: `fiber/v2` → `fiber/v3`
- Todas as assinaturas de método: `*fiber.Ctx` → `fiber.Ctx`
- `ctx.BodyParser(&createUserDTO)` → `ctx.Bind().Body(&createUserDTO)`
- `ctx.QueryParser(&userFiltersDTO)` → `ctx.Bind().Query(&userFiltersDTO)`
- `ctx.Params("userID")` — permanece igual

---

#### T2.8 — Atualizar `internal/infra/entrypoint/rest/group_controller.go`

**Mudanças**:
- Import: `fiber/v2` → `fiber/v3`; adicionar import `jwtware "github.com/gofiber/contrib/v3/jwt"`
- Todas as assinaturas de método: `*fiber.Ctx` → `fiber.Ctx`
- `ctx.BodyParser(...)` → `ctx.Bind().Body(...)`
- `ctx.QueryParser(...)` → `ctx.Bind().Query(...)`
- `ctx.Locals("user")` → `jwtware.FromContext(ctx)` (em todos os métodos que chamam `GetAuthUserID`)
- `ctx.Params(...)` — permanece igual

---

#### T2.9 — Atualizar `internal/infra/entrypoint/rest/group_invite_controller.go`

**Mudanças**:
- Import: `fiber/v2` → `fiber/v3`; adicionar import `jwtware "github.com/gofiber/contrib/v3/jwt"`
- Todas as assinaturas de método: `*fiber.Ctx` → `fiber.Ctx`
- `ctx.Locals("user")` → `jwtware.FromContext(ctx)` (em todos os métodos que chamam `GetAuthUserID`)
- `ctx.Params(...)` — permanece igual

---

#### T2.10 — Atualizar arquivos de teste

**Arquivos impactados**:
- `internal/infra/entrypoint/error_handler_test.go`
- `internal/infra/entrypoint/rest/auth_controller_test.go`
- `internal/infra/entrypoint/rest/user_controller_test.go`
- `internal/infra/entrypoint/rest/group_controller_test.go`
- `internal/infra/entrypoint/rest/group_invite_controller_test.go`

**Mudanças em cada arquivo**:
- Import: `fiber/v2` → `fiber/v3`
- Handler lambdas inline: `func(c *fiber.Ctx) error` → `func(c fiber.Ctx) error`
- `fiber.New(fiber.Config{...})` — permanece igual
- Chamadas a `httptest.NewRequest` e `app.Test` — verificar se assinatura de `Test` mudou em V3

> **Nota**: Em Fiber V3, `app.Test(req, config ...fiber.TestConfig)` substitui `app.Test(req, timeout ...time.Duration)`. Se os testes usarem o segundo argumento de timeout, precisará ser adaptado para `fiber.TestConfig{Timeout: ...}`.

---

#### T2.11 — Verificar build, testes e documentação

```bash
make build
make test
make generate-docs
```

Todos devem passar. Corrigir qualquer falha de compilação antes de avançar.

---

#### T2.12 — Atualizar a Constituição do projeto

Arquivo: `.specify/memory/constitution.md`

**Mudanças necessárias**:

| Localização | Texto atual | Novo texto |
|-------------|-------------|------------|
| Princípio IV — assinatura de handler | `func (c *XController) Method(ctx *fiber.Ctx) error` | `func (c *XController) Method(ctx fiber.Ctx) error` |
| Princípio IV — fluxo de requisição | `BodyParser →` | `Bind().Body() →` |
| Princípio IV — falha de parsing | `Falha em BodyParser DEVE retornar fiber.NewError(fiber.StatusUnprocessableEntity)` | `Falha em Bind().Body() DEVE retornar fiber.NewError(fiber.StatusUnprocessableEntity)` |
| Princípio IV — extração de userID | `c.AuthTokenManager.GetAuthUserID(ctx.Locals("user"))` | `c.AuthTokenManager.GetAuthUserID(jwtware.FromContext(ctx))` |
| Padrões Tecnológicos | `Fiber v2 — usar apenas padrões idiomáticos do Fiber` | `Fiber v3 — usar apenas padrões idiomáticos do Fiber` |

Incrementar `CONSTITUTION_VERSION` de `1.1.0` para `1.2.0` (novo padrão de handler adicionado — MINOR) e atualizar `LAST_AMENDED_DATE` para `2026-06-21`.
