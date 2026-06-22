# Implementation Plan: Autenticação Dual-Channel com Cookie HttpOnly

**Branch**: `004-jwt-cookie-auth` | **Date**: 2026-06-22 | **Spec**: [spec.md](spec.md)

## Resumo

Adicionar suporte a autenticação JWT via cookie httpOnly como alternativa ao cabeçalho `Authorization`, eliminando a exposição do token no localStorage do frontend. O middleware de autenticação passa a aceitar cookie (com precedência) ou cabeçalho, mantendo retrocompatibilidade total com clientes mobile. Dois novos endpoints são adicionados: `POST /logout` (remove o cookie) e `GET /users/me` (retorna dados do usuário autenticado).

## Technical Context

**Language/Version**: Go 1.26.4  
**Primary Dependencies**: Fiber v3 (v3.3.0), gofiber/contrib/v3/jwt (v1.1.6), fiber/v3/extractors (transitivo), golang-jwt/jwt v5.3.1  
**Storage**: PostgreSQL — sem migrações (feature é puramente HTTP)  
**Testing**: testify/assert, go.uber.org/mock/mockgen  
**Target Platform**: Linux server (Docker)  
**Performance Goals**: p95 < 200ms (meta existente por constituição)  
**Constraints**: Cookie httpOnly, Secure configurável via env, SameSite=Lax (cross-origin entre subdomínios)  
**Scale/Scope**: Afeta todos os endpoints protegidos existentes (mudança de middleware)

## Constitution Check

| Princípio | Status | Observação |
|-----------|--------|------------|
| I. Clean Architecture | ✅ Pass | Cookie management em `infra/entrypoint`; auth business logic permanece em `application/domain` |
| II. Disciplina de Testes | ✅ Pass | Testes unitários obrigatórios por handler e para o extractor do middleware |
| III. Validação no Domínio | ✅ Pass | Nenhuma alteração em entidades de domínio; `AuthSession` e `User` inalterados |
| IV. Contrato de API Consistente | ✅ Pass | Handlers seguem o fluxo padrão; anotações Swagger atualizadas |
| V. Abstração de Infraestrutura | ✅ Pass | Cookie ops em `entrypoint` (camada HTTP); `AuthTokenManager` injetado onde necessário |
| VI. Simplicidade & YAGNI | ✅ Pass | `extractors.Chain` elimina código custom; sem novas abstrações além do necessário |
| VII. Performance & Observabilidade | ✅ Pass | `context.Context` mantido em todos os métodos; logging de erros preservado |
| VIII. Idioma dos Artefatos | ✅ Pass | Artefatos em pt-BR; código em inglês |

**Violations**: Nenhuma

## Project Structure

### Documentação (esta feature)

```text
specs/004-jwt-cookie-auth/
├── plan.md              # Este arquivo
├── research.md          # Pesquisa sobre jwtware extractors e política de cookie
├── data-model.md        # Sem alterações de schema; novo campo de env
├── quickstart.md        # Guia de verificação local
├── contracts/
│   └── api.md           # Contratos dos 3 endpoints (login atualizado, logout, /users/me)
└── tasks.md             # Gerado por /speckit.tasks
```

### Source Code (alterações por arquivo)

```text
internal/
├── infra/
│   ├── config/
│   │   └── config.go                           # + CookieSecure bool em AuthConfig
│   ├── entrypoint/
│   │   ├── middlewares.go                      # NewAuthMiddleware com extractors.Chain
│   │   ├── cookie.go                           # [NOVO] constante authCookieName + helpers setCookie/clearCookie
│   │   └── rest/
│   │       ├── auth_controller.go              # Login seta cookie; + handler Logout
│   │       ├── auth_controller_test.go         # + cenários de cookie e logout
│   │       ├── user_controller.go              # + handler GetMe; injetar AuthTokenManager
│   │       └── user_controller_test.go         # + cenários GetMe
│   └── runner.go                               # Passar jwtAuthTokenManager e cookieSecure
└── (domain/ e application/ — sem alterações)

.env.example                                    # + AUTH_COOKIE_SECURE=true
```

**Structure Decision**: Single project, camada `infra/entrypoint` exclusivamente. Nenhum novo pacote criado.

---

## Checkpoints

### Checkpoint 1 (CP1): Cookie na Autenticação e Middleware Dual-Channel

**Escopo**: User Stories P1 e P4 — login define cookie seguro; todas as rotas protegidas aceitam cookie ou cabeçalho com cookie tendo precedência.

**Arquivos alterados**:

#### `internal/infra/config/config.go`
- Adicionar `CookieSecure bool` em `AuthConfig` com tag `env:"AUTH_COOKIE_SECURE" envDefault:"true"`

#### `.env.example`
- Adicionar `AUTH_COOKIE_SECURE=true`

#### `internal/infra/entrypoint/cookie.go` (novo arquivo)
- Constante `authCookieName = "access_token"`
- Função `setCookie(ctx fiber.Ctx, token string, expiresIn int64, secure bool)` — define o cookie httpOnly
- Função `clearCookie(ctx fiber.Ctx)` — remove o cookie (Expires no passado, MaxAge=-1)

#### `internal/infra/entrypoint/middlewares.go`
- Atualizar `NewAuthMiddleware` para usar `jwtware.Config{Extractor: extractors.Chain(extractors.FromCookie(authCookieName), extractors.FromAuthHeader("Bearer")), ...}`
- Import de `"github.com/gofiber/fiber/v3/extractors"`
- Verificar que o `ErrorHandler` existente mapeia `extractors.ErrNotFound` para 400 (já está correto via `jwtware.ErrMissingToken`)

#### `internal/infra/entrypoint/rest/auth_controller.go`
- Adicionar campo `cookieSecure bool` na struct `AuthController`
- Atualizar `NewAuthController(authService application.AuthService, cookieSecure bool) *AuthController`
- Em `Login`: após `authSessionDTO`, chamar `setCookie(ctx, authSession.AccessToken, authSession.ExpiresIn, c.cookieSecure)` antes de retornar

#### `internal/infra/entrypoint/rest/auth_controller_test.go`
- Atualizar chamadas a `NewAuthController` para passar `cookieSecure`
- Adicionar cenário: `should set auth cookie on successful login`
- Verificar header `Set-Cookie` na resposta com atributos corretos

#### `internal/infra/runner.go`
- Passar `cfg.Auth.CookieSecure` para `NewAuthController`

**Testes de integração do middleware** (em `auth_controller_test.go` ou arquivo separado `middlewares_test.go`):
- `should authenticate via cookie when cookie is present`
- `should authenticate via header when cookie is absent`
- `should reject when cookie is present but invalid (no fallback to header)`
- `should reject when neither cookie nor header present`

---

### Checkpoint 2 (CP2): Endpoint de Logout

**Escopo**: User Story P2 — logout remove o cookie do navegador.

**Arquivos alterados**:

#### `internal/infra/entrypoint/rest/auth_controller.go`
- Adicionar handler `Logout(ctx fiber.Ctx) error`:
  - Chamar `clearCookie(ctx)`
  - Retornar `ctx.SendStatus(fiber.StatusNoContent)`

#### `internal/infra/entrypoint/rest/auth_controller_test.go`
- `Test_AuthController_Logout`:
  - `should clear auth cookie and return 204 on logout`
  - Verificar header `Set-Cookie` com `MaxAge=-1` e `Expires` no passado

#### `internal/infra/entrypoint/routes.go`
- Registrar `api.Post("/logout", authController.Logout)` **antes** de `api.Use(authMiddleware)` (logout não requer autenticação)

---

### Checkpoint 3 (CP3): Endpoint GET /users/me

**Escopo**: User Story P3 — usuário consulta seus próprios dados sem informar ID.

**Arquivos alterados**:

#### `internal/infra/entrypoint/rest/user_controller.go`
- Adicionar campo `authTokenManager domain.AuthTokenManager` na struct `UserController`
- Atualizar `NewUserController` para receber `authTokenManager domain.AuthTokenManager` como parâmetro adicional
- Adicionar handler `GetMe(ctx fiber.Ctx) error`:
  ```
  authUserID ← authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))
  user ← userService.GetByID(ctx.Context(), authUserID)
  userDTO ← mapUserFromDomain(user)
  return ctx.JSON(userDTO)
  ```

#### `internal/infra/entrypoint/rest/user_controller_test.go`
- `Test_UserController_GetMe`:
  - `should return authenticated user data successfully`
  - `should return unauthorized when token extraction fails`
  - `should return not_found when user does not exist`

#### `internal/infra/entrypoint/routes.go`
- Registrar `api.Get("/users/me", userController.GetMe)` **antes** de `api.Get("/users/:userID", userController.GetByID)` (ambos após o authMiddleware)

#### `internal/infra/runner.go`
- Passar `jwtAuthTokenManager` para `NewUserController`

#### Swagger (`routes.go`)
- Adicionar anotação para `GET /api/v1/users/me`
- Adicionar anotação para `POST /api/v1/logout`
- Atualizar anotação de `POST /api/v1/login` para documentar o Set-Cookie na resposta

---

## Gates de Qualidade (por checkpoint)

Executar após cada checkpoint antes de abrir PR:

1. `make build` — compila sem erros
2. `make test` — todos os testes passam
3. `make generate-docs` — spec Swagger gerada sem erros
4. Nenhum placeholder de template em spec ou plano
