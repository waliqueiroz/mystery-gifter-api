# Tasks: Autenticação Dual-Channel com Cookie HttpOnly

**Input**: Design documents from `/specs/004-jwt-cookie-auth/`
**Prerequisites**: plan.md ✅ spec.md ✅ research.md ✅ data-model.md ✅ contracts/api.md ✅

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Pode ser executado em paralelo (arquivos diferentes, sem dependências pendentes)
- **[Story]**: User story à qual a tarefa pertence (US1=Login Seguro, US2=Logout, US3=Perfil Logado, US4=Mobile)

---

## Phase 1: Setup (Configuração Inicial)

**Objetivo**: Preparar configuração de ambiente antes de qualquer implementação

- [x] T001 Adicionar campo `CookieSecure bool` com tag `env:"AUTH_COOKIE_SECURE" envDefault:"true"` à struct `AuthConfig` em `internal/infra/config/config.go`
- [x] T002 [P] Adicionar linha `AUTH_COOKIE_SECURE=true` ao arquivo `.env.example`

---

## Phase 2: Foundational (Infraestrutura de Cookie e Middleware)

**Objetivo**: Criar os blocos de construção que todas as user stories dependem

**⚠️ CRÍTICO**: Nenhuma user story pode ser implementada antes desta fase estar completa

- [x] T003 Criar arquivo `internal/infra/entrypoint/cookie.go` com constante `authCookieName = "access_token"`, função `setCookie(ctx fiber.Ctx, token string, expiresIn int64, secure bool)` definindo cookie httpOnly/SameSite=Lax, e função `clearCookie(ctx fiber.Ctx)` expirando o cookie com `MaxAge: -1`
- [x] T004 Atualizar `NewAuthMiddleware` em `internal/infra/entrypoint/middlewares.go` adicionando `Extractor: extractors.Chain(extractors.FromCookie(authCookieName), extractors.FromAuthHeader("Bearer"))` à `jwtware.Config` e importando `"github.com/gofiber/fiber/v3/extractors"`

**Checkpoint**: Cookie helpers criados e middleware configurado com dual-channel — user stories podem ser iniciadas

---

## Phase 3: User Story 1 + User Story 4 — Login com Cookie Seguro e Retrocompatibilidade (P1 e P4) 🎯 MVP

**Goal**: Login define cookie httpOnly seguro no navegador; todas as rotas protegidas aceitam cookie (com precedência) ou cabeçalho `Authorization: Bearer`

**Independent Test**: Realizar login e verificar header `Set-Cookie` na resposta com atributos `HttpOnly` e `SameSite=Lax`; confirmar que acessar rota protegida via cookie retorna 200; confirmar que acessar via `Authorization: Bearer` também retorna 200

### Implementação para US1 e US4

- [x] T005 [US1] Adicionar campo `cookieSecure bool` à struct `AuthController` e atualizar assinatura `NewAuthController(authService application.AuthService, cookieSecure bool) *AuthController` em `internal/infra/entrypoint/rest/auth_controller.go`
- [x] T006 [US1] Atualizar handler `Login` em `internal/infra/entrypoint/rest/auth_controller.go` para chamar `setCookie(ctx, authSession.AccessToken, authSession.ExpiresIn, c.cookieSecure)` logo antes de `return ctx.JSON(authSessionDTO)`
- [x] T007 [US1] Atualizar `NewAuthController` em `internal/infra/runner.go` passando `cfg.Auth.CookieSecure` como segundo argumento
- [x] T008 [P] [US1] Atualizar `internal/infra/entrypoint/rest/auth_controller_test.go`: corrigir todas as chamadas a `NewAuthController` para o novo número de argumentos e adicionar cenário `"should set auth cookie on successful login"` verificando presença do header `Set-Cookie` com nome `access_token`, atributos `HttpOnly` e `SameSite=Lax`
- [x] T009 [P] [US1/US4] Criar arquivo `internal/infra/entrypoint/middlewares_test.go` com `Test_NewAuthMiddleware` cobrindo: autenticação via cookie válido (200), autenticação via header Bearer válido sem cookie (200), cookie presente mas com token inválido sem fallback ao header (401), ausência de cookie e header (400)
- [x] T010 [US1] Executar `make test` e corrigir todas as falhas antes de avançar

**Checkpoint**: Login define cookie seguro; rotas protegidas aceitam cookie e header; retrocompatibilidade mobile validada ✓

---

## Phase 4: User Story 2 — Encerramento de Sessão Explícito (P2)

**Goal**: `POST /logout` remove o cookie de autenticação do navegador; não requer autenticação

**Independent Test**: Chamar `POST /logout` e verificar resposta 204 com header `Set-Cookie: access_token=; MaxAge=-1` indicando remoção do cookie

### Implementação para US2

- [x] T011 [US2] Adicionar handler `Logout(ctx fiber.Ctx) error` em `internal/infra/entrypoint/rest/auth_controller.go` que chama `clearCookie(ctx)` e retorna `ctx.SendStatus(fiber.StatusNoContent)`
- [x] T012 [US2] Registrar rota `api.Post("/logout", authController.Logout)` em `internal/infra/entrypoint/routes.go` **antes** de `api.Use(authMiddleware)` (logout não requer autenticação)
- [x] T013 [US2] Adicionar anotação Swagger para `POST /api/v1/logout` em `internal/infra/entrypoint/routes.go` com tag `auth`, resposta 204 e nota sobre remoção do cookie
- [x] T014 [P] [US2] Adicionar função `Test_AuthController_Logout` em `internal/infra/entrypoint/rest/auth_controller_test.go` com cenários: `"should clear auth cookie and return 204 on logout"` verificando status 204 e header `Set-Cookie` com `MaxAge=-1`
- [x] T015 [US2] Executar `make test` e `make generate-docs` e corrigir eventuais falhas

**Checkpoint**: Logout funcional; cookie removido na resposta; endpoint documentado no Swagger ✓

---

## Phase 5: User Story 3 — Consulta de Dados do Usuário Logado (P3)

**Goal**: `GET /users/me` retorna `UserDTO` do usuário autenticado sem exigir ID na URL; deve ser declarado antes de `GET /users/:userID`

**Independent Test**: Chamar `GET /users/me` com autenticação válida (cookie ou header) e verificar resposta 200 com os campos `id`, `name`, `surname`, `email`, `created_at`, `updated_at`

### Implementação para US3

- [x] T016 [US3] Adicionar campo `authTokenManager domain.AuthTokenManager` à struct `UserController` e atualizar `NewUserController` em `internal/infra/entrypoint/rest/user_controller.go` para receber `authTokenManager domain.AuthTokenManager` como terceiro parâmetro após `bcryptPasswordManager`
- [x] T017 [US3] Adicionar handler `GetMe(ctx fiber.Ctx) error` em `internal/infra/entrypoint/rest/user_controller.go`: extrair `userID` via `c.authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))`; chamar `c.userService.GetByID(ctx.Context(), userID)`; mapear com `mapUserFromDomain` e retornar `ctx.JSON(userDTO)`
- [x] T018 [US3] Registrar rota `api.Get("/users/me", userController.GetMe)` em `internal/infra/entrypoint/routes.go` **antes** da linha `api.Get("/users/:userID", ...)` (ambas após `api.Use(authMiddleware)`)
- [x] T019 [US3] Adicionar anotação Swagger para `GET /api/v1/users/me` em `internal/infra/entrypoint/routes.go` com tag `users`, security `Bearer`, resposta 200 referenciando `UserDTO`, 401 e 404
- [x] T020 [US3] Atualizar `NewUserController` em `internal/infra/runner.go` passando `jwtAuthTokenManager` como terceiro argumento
- [x] T021 [P] [US3] Adicionar função `Test_UserController_GetMe` em `internal/infra/entrypoint/rest/user_controller_test.go` com cenários: `"should return authenticated user data successfully"` (200 com UserDTO correto), `"should return unauthorized when token extraction fails"` (401), `"should return not_found when user does not exist"` (404); atualizar chamadas existentes a `NewUserController` para o novo número de argumentos
- [x] T022 [US3] Executar `make test` e `make generate-docs` e corrigir eventuais falhas

**Checkpoint**: `GET /users/me` funcional para cookie e header; sem conflito com `GET /users/:userID`; documentado no Swagger ✓

---

## Phase 6: Polish & Validação Final

**Objetivo**: Garantir consistência, documentação e validação end-to-end da feature completa

- [x] T023 Atualizar anotação Swagger de `POST /api/v1/login` em `internal/infra/entrypoint/routes.go` adicionando nota na `description` sobre o cookie `access_token` definido na resposta (`HttpOnly`, `SameSite=Lax`)
- [x] T024 [P] Executar `make build` para confirmar compilação limpa
- [x] T025 [P] Executar `make test` final para confirmar que todos os testes passam
- [x] T026 Executar `make generate-docs` e `make serve-docs` e validar os novos endpoints `POST /logout` e `GET /users/me` na UI do Swagger em `http://localhost:8081`
- [ ] T027 Seguir o roteiro de `specs/004-jwt-cookie-auth/quickstart.md` para validar o fluxo completo: login → cookie definido → acesso com cookie → logout → acesso negado

---

## Dependências e Ordem de Execução

### Dependências entre Fases

- **Phase 1 (Setup)**: Sem dependências — pode começar imediatamente
- **Phase 2 (Foundational)**: Depende de Phase 1 — **bloqueia todas as user stories**
- **Phase 3 (US1+US4)**: Depende de Phase 2
- **Phase 4 (US2)**: Depende de Phase 2 (reutiliza `clearCookie` de T003); pode ser paralela a Phase 3
- **Phase 5 (US3)**: Depende de Phase 2; pode ser paralela a Phase 3 e 4
- **Phase 6 (Polish)**: Depende da conclusão de todas as fases anteriores

### Dependências entre User Stories

- **US1 (P1)**: Depende de Phase 2; sem dependências de outras user stories
- **US2 (P2)**: Depende de Phase 2 (usa `clearCookie`); sem dependências de US1/US3
- **US3 (P3)**: Depende de Phase 2; sem dependências de US1/US2
- **US4 (P4)**: Coberta pela Phase 2 (middleware) e pela Phase 3 (testes)

### Dentro de Cada User Story

- Implementação antes dos testes (ou paralela)
- Runner update (`runner.go`) sempre após mudança de construtor
- `make test` obrigatório antes de avançar para próximo checkpoint
- `make generate-docs` após qualquer mudança de rota ou DTO

### Oportunidades de Paralelismo

- T001 e T002 (Phase 1) podem ser executadas em paralelo
- T008 e T009 (Phase 3, arquivos diferentes) podem ser executadas em paralelo
- Após Phase 2 completa: Phase 3, 4 e 5 podem ser iniciadas em paralelo por agentes/desenvolvedores diferentes
- T024 e T025 (Phase 6) podem ser executadas em paralelo

---

## Exemplo de Paralelismo: Phase 3 (US1 + US4)

```
# Após T005, T006, T007 (sequenciais):

Agente A: T008 — Testes de auth_controller (auth_controller_test.go)
Agente B: T009 — Testes de middleware   (middlewares_test.go)

# Aguardar ambos completarem → T010 (make test)
```

---

## Estratégia de Implementação

### MVP (apenas US1 + US4 — Cookie no Login)

1. Completar Phase 1 (Setup)
2. Completar Phase 2 (Foundational)
3. Completar Phase 3 (US1 + US4)
4. **PARAR e VALIDAR**: Cookie definido no login, rotas protegidas funcionam com cookie e header
5. Entregar CP1 via PR

### Entrega Incremental

1. Setup + Foundational → infraestrutura base pronta
2. US1+US4 (Phase 3) → MVP: login seguro + retrocompatibilidade mobile ✅
3. US2 (Phase 4) → logout funcional ✅
4. US3 (Phase 5) → endpoint `/users/me` ✅
5. Polish (Phase 6) → validação completa ✅

---

## Notas

- `[P]` = arquivos diferentes, sem dependências pendentes entre si
- `[USn]` = user story correspondente para rastreabilidade
- `make test` é obrigatório ao final de cada phase — não avançar com testes falhando
- `make generate-docs` obrigatório após qualquer adição/mudança de rota ou DTO Swagger
- A rota `GET /users/me` DEVE ser registrada antes de `GET /users/:userID` — esta ordem é crítica
- A rota `POST /logout` DEVE ser registrada antes de `api.Use(authMiddleware)` — esta ordem é crítica
