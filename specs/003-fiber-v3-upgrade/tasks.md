# Tasks: Atualização Go 1.26.4 + Fiber V2 → V3

**Input**: `/specs/003-fiber-v3-upgrade/` (plan.md, spec.md, research.md)  
**Branch**: `003-fiber-v3-upgrade`

## Formato: `[ID] [P?] [Story?] Descrição com caminho do arquivo`

- **[P]**: Executável em paralelo (arquivos distintos, sem dependência de task incompleta)
- **[Story]**: User story correspondente da spec.md

---

## Phase 1: User Story 1 — Atualização do Go para 1.26.4 (Priority: P1) 🎯

**Goal**: Atualizar o toolchain do projeto para Go 1.26.4, isolando essa mudança das demais.

**Independent Test**: `make build && make test` passam sem erros com Go 1.26.4.

**PR boundary**: Este checkpoint vira um PR independente antes de iniciar o Checkpoint 2.

- [ ] T001 [US1] Atualizar diretiva `go 1.25.1` para `go 1.26.4` em `go.mod`
- [ ] T002 [US1] Executar `go mod tidy` para sincronizar `go.sum` com o novo toolchain; atualizar dependências indiretas que falharem a compilação
- [ ] T003 [US1] Verificar build e testes: `make build && make test` — ambos devem passar sem erros

**Checkpoint 1**: Go 1.26.4 funcionando. Abrir PR, aguardar merge antes de iniciar Phase 2.

---

## Phase 2: User Stories 2, 3 e 4 — Migração Fiber V2 → V3 (Priority: P2–P4)

**Goal**: Atualizar todas as dependências e adaptar o código de infra para o Fiber V3, mantendo comportamento idêntico da API.

**Independent Test (US2)**: Todos os endpoints retornam os mesmos status HTTP e estruturas de resposta.  
**Independent Test (US3)**: Requisições com token JWT válido, inválido e ausente retornam os status esperados.  
**Independent Test (US4)**: `make test` passa — o middleware recover continua capturando panics.

**PR boundary**: Este checkpoint vira um PR stacked no Checkpoint 1.

### Atualização de dependências

- [ ] T004 [US2] Atualizar dependências Fiber em `go.mod`: `go get github.com/gofiber/fiber/v3@latest && go get github.com/gofiber/contrib/v3/jwt@latest && go mod tidy`

### Camada de infraestrutura (paralelos após T004)

- [ ] T005 [P] [US2] Atualizar `internal/infra/runner.go`: imports `fiber/v2` → `fiber/v3`; `fiber/v2/middleware/cors` → `fiber/v3/middleware/cors`; `fiber/v2/middleware/recover` → `fiber/v3/middleware/recover`
- [ ] T006 [P] [US3] Atualizar `internal/infra/entrypoint/middlewares.go`: import `gofiber/contrib/jwt` → `gofiber/contrib/v3/jwt`; import `fiber/v2` → `fiber/v3`; assinatura do `ErrorHandler` de `func(c *fiber.Ctx, err error) error` para `func(c fiber.Ctx, err error) error`
- [ ] T007 [P] [US2] Atualizar `internal/infra/entrypoint/error_handler.go`: import `fiber/v2` → `fiber/v3`; assinaturas de `CustomErrorHandler` e `sendError` de `*fiber.Ctx` para `fiber.Ctx`
- [ ] T008 [P] [US2] Atualizar `internal/infra/entrypoint/routes.go`: import `fiber/v2` → `fiber/v3`

### Controllers (paralelos entre si, após T004)

- [ ] T009 [P] [US2] Atualizar `internal/infra/entrypoint/rest/auth_controller.go`: import `fiber/v2` → `fiber/v3`; assinatura de `Login` de `*fiber.Ctx` para `fiber.Ctx`; `ctx.BodyParser(&credentialsDTO)` → `ctx.Bind().Body(&credentialsDTO)`
- [ ] T010 [P] [US2] Atualizar `internal/infra/entrypoint/rest/user_controller.go`: import `fiber/v2` → `fiber/v3`; todas as assinaturas de método de `*fiber.Ctx` para `fiber.Ctx`; `ctx.BodyParser(...)` → `ctx.Bind().Body(...)`; `ctx.QueryParser(...)` → `ctx.Bind().Query(...)`
- [ ] T011 [P] [US3] Atualizar `internal/infra/entrypoint/rest/group_controller.go`: import `fiber/v2` → `fiber/v3`; adicionar import `jwtware "github.com/gofiber/contrib/v3/jwt"`; todas as assinaturas de método de `*fiber.Ctx` para `fiber.Ctx`; `ctx.BodyParser(...)` → `ctx.Bind().Body(...)`; `ctx.QueryParser(...)` → `ctx.Bind().Query(...)`; todas as ocorrências de `ctx.Locals("user")` → `jwtware.FromContext(ctx)`
- [ ] T012 [P] [US3] Atualizar `internal/infra/entrypoint/rest/group_invite_controller.go`: import `fiber/v2` → `fiber/v3`; adicionar import `jwtware "github.com/gofiber/contrib/v3/jwt"`; todas as assinaturas de método de `*fiber.Ctx` para `fiber.Ctx`; todas as ocorrências de `ctx.Locals("user")` → `jwtware.FromContext(ctx)`

### Arquivos de teste (paralelos entre si, após T004)

- [ ] T013 [P] [US2] Atualizar `internal/infra/entrypoint/error_handler_test.go`: import `fiber/v2` → `fiber/v3`; handler lambdas de `func(c *fiber.Ctx) error` para `func(c fiber.Ctx) error`
- [ ] T014 [P] [US2] Atualizar `internal/infra/entrypoint/rest/auth_controller_test.go`: import `fiber/v2` → `fiber/v3`; handler lambdas de `*fiber.Ctx` para `fiber.Ctx`
- [ ] T015 [P] [US2] Atualizar `internal/infra/entrypoint/rest/user_controller_test.go`: import `fiber/v2` → `fiber/v3`; handler lambdas de `*fiber.Ctx` para `fiber.Ctx`
- [ ] T016 [P] [US3] Atualizar `internal/infra/entrypoint/rest/group_controller_test.go`: import `fiber/v2` → `fiber/v3`; handler lambdas de `*fiber.Ctx` para `fiber.Ctx`
- [ ] T017 [P] [US3] Atualizar `internal/infra/entrypoint/rest/group_invite_controller_test.go`: import `fiber/v2` → `fiber/v3`; handler lambdas de `*fiber.Ctx` para `fiber.Ctx`

### Verificação e encerramento

- [ ] T018 [US2] Verificar build e testes: `make build && make test` — ambos devem passar sem erros
- [ ] T019 [US2] Verificar geração de documentação: `make generate-docs` — deve concluir sem erros

**Checkpoint 2**: Migração Fiber V3 completa e validada.

---

## Phase 3: Polish — Atualização da Constituição do Projeto

**Goal**: Manter a Constituição como fonte canônica de padrões, refletindo os padrões Fiber V3.

- [ ] T020 Atualizar `.specify/memory/constitution.md`: substituir `*fiber.Ctx` por `fiber.Ctx` na assinatura de handler (Princípio IV); substituir `BodyParser →` por `Bind().Body() →` no fluxo de requisição (Princípio IV); substituir `Falha em BodyParser` por `Falha em Bind().Body()` (Princípio IV); substituir `ctx.Locals("user")` por `jwtware.FromContext(ctx)` na extração de userID (Princípio IV); substituir `Fiber v2` por `Fiber v3` nos Padrões Tecnológicos; incrementar `CONSTITUTION_VERSION` de `1.1.0` para `1.2.0` e atualizar `LAST_AMENDED_DATE` para `2026-06-21`

---

## Dependências e Ordem de Execução

### Dependências entre fases

- **Phase 1 (Checkpoint 1)**: Sem dependências — iniciar imediatamente
- **Phase 2 (Checkpoint 2)**: Requer Checkpoint 1 mergeado — branch stacked no branch do Checkpoint 1
- **Phase 3 (Polish)**: Requer Phase 2 completa

### Dependências internas do Checkpoint 2

```
T004 (deps)
  ├── T005 [P] (runner.go)
  ├── T006 [P] (middlewares.go)
  ├── T007 [P] (error_handler.go)
  ├── T008 [P] (routes.go)
  ├── T009 [P] (auth_controller.go)
  ├── T010 [P] (user_controller.go)
  ├── T011 [P] (group_controller.go)
  ├── T012 [P] (group_invite_controller.go)
  ├── T013 [P] (error_handler_test.go)
  ├── T014 [P] (auth_controller_test.go)
  ├── T015 [P] (user_controller_test.go)
  ├── T016 [P] (group_controller_test.go)
  └── T017 [P] (group_invite_controller_test.go)
       └── T018 (make build && make test)
            └── T019 (make generate-docs)
```

### Oportunidades de paralelismo — Checkpoint 2

Após T004 (atualização de deps), as tasks T005–T017 podem todas ser executadas em paralelo entre si (arquivos distintos, sem dependência mútua):

```bash
# Paralelo — infraestrutura (após T004)
Task: "T005 — runner.go"
Task: "T006 — middlewares.go"
Task: "T007 — error_handler.go"
Task: "T008 — routes.go"

# Paralelo — controllers (após T004)
Task: "T009 — auth_controller.go"
Task: "T010 — user_controller.go"
Task: "T011 — group_controller.go"
Task: "T012 — group_invite_controller.go"

# Paralelo — testes (após T004)
Task: "T013 — error_handler_test.go"
Task: "T014 — auth_controller_test.go"
Task: "T015 — user_controller_test.go"
Task: "T016 — group_controller_test.go"
Task: "T017 — group_invite_controller_test.go"
```

---

## Estratégia de Implementação

### Entrega sequencial por checkpoint (stacked PRs)

1. **Checkpoint 1 — Go 1.26.4** (T001–T003): branch em `main`, PR independente
2. **Checkpoint 2 — Fiber V3** (T004–T019): branch stacked no Checkpoint 1
3. **Polish** (T020): incluído no PR do Checkpoint 2 ou PR separado em `main` após merges

### Validação em cada checkpoint

- **Após Checkpoint 1**: `make build && make test` passam → merge do PR
- **Após Checkpoint 2**: `make build && make test && make generate-docs` passam → merge do PR
- **Após Polish**: Constituição reflete Fiber V3 → feature concluída

---

## Notas

- `[P]` = arquivos distintos, sem dependência de task incompleta
- `app.Test(req)` é chamado sem segundo argumento em todos os testes — sem breaking change na assinatura V3 de `app.Test`
- `cors.New()` e `recover.New()` usam config padrão — nenhuma mudança de comportamento
- `AuthTokenManager.GetAuthUserID(token any)` não muda — `*jwt.Token` satisfaz `any`
- Quaisquer outras dependências (`sqlx`, `squirrel`, etc.) só serão atualizadas se falharem compilação
