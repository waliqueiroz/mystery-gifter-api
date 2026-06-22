# Research: Autenticação Dual-Channel com Cookie HttpOnly

## 1. Extração de Token via Cookie no Fiber v3

**Decisão**: Usar `extractors.Chain(extractors.FromCookie("access_token"), extractors.FromAuthHeader("Bearer"))` no campo `Extractor` do `jwtware.Config`.

**Rationale**: O pacote `github.com/gofiber/fiber/v3/extractors` (já transitivamente disponível via `gofiber/contrib/v3/jwt`) expõe exatamente as primitivas necessárias:
- `extractors.FromCookie(key)` — lê o cookie pelo nome; retorna `ErrNotFound` se ausente
- `extractors.FromAuthHeader(scheme)` — lê o cabeçalho `Authorization: Bearer <token>`
- `extractors.Chain(e1, e2, ...)` — tenta cada extractor em ordem, retorna o primeiro com sucesso

Quando o cookie está presente mas contém um JWT inválido, `FromCookie` retorna o valor corrompido (sem erros de extração), e o jwtware rejeita com 401 durante a validação criptográfica — sem fazer fallback para o cabeçalho. Este é exatamente o comportamento especificado em FR-005 e no edge case "cookie inválido + header válido".

**Alternativas consideradas**:
- `Config.TokenProcessorFunc`: transforma o token após extração, não muda a fonte — descartado
- Implementar um extractor custom com `extractors.FromCustom`: funcionaria, mas `Chain` já resolve o problema sem código adicional — descartado por YAGNI

---

## 2. Configuração do Cookie de Sessão

**Decisão**: Usar `fiber.Cookie` com os seguintes atributos:
- `Name: "access_token"` — consistente com o campo `access_token` do `AuthSessionDTO`
- `Value: authSession.AccessToken`
- `Expires: time.Unix(authSession.ExpiresIn, 0)` — aproveita o campo já disponível no `AuthSession`
- `HTTPOnly: true` — impede acesso via JavaScript (FR-002)
- `Secure: <configurável via env>` — `true` em produção, `false` em desenvolvimento (FR-003)
- `SameSite: "Lax"` — permite requisições cross-site de mesmo domínio pai (ex: `app.x.com` → `api.x.com`), bloqueando cross-site POST de domínios externos (FR-004)

**Rationale para SameSite=Lax**: Frontend e API operam em origens distintas mas compartilham o mesmo domínio pai (subdomínios). `SameSite=Lax` permite que cookies sejam enviados neste cenário enquanto protege contra CSRF de domínios externos. `SameSite=Strict` bloquearia navegações legítimas entre subdomínios. `SameSite=None` exigiria `Secure=true` sempre e seria menos restritivo que o necessário.

**Alternativas consideradas**:
- `SameSite=Strict`: bloqueia cookies em navegações cross-site mesmo dentro do mesmo domínio pai — mais restritivo do que necessário para o cenário declarado — descartado
- `SameSite=None`: requer `Secure=true` e permite qualquer origem cross-site — menos seguro — descartado

---

## 3. Gerenciamento do Cookie no Logout

**Decisão**: O endpoint de logout limpa o cookie usando `ctx.Cookie` com `Expires` no passado (convenção HTTP para remoção de cookie). Retorna HTTP 204 No Content. **Não requer autenticação** — é registrado antes do `authMiddleware` nas rotas.

**Rationale**: Como não há revogação server-side (decisão documentada na spec), o logout é puramente uma operação de limpeza de cookie no cliente. Exigir autenticação para logout criaria um cenário ruim: usuário com token expirado não consegue fazer logout explícito (limpar o cookie). Registrar antes do middleware evita este problema.

**Alternativas consideradas**:
- Logout atrás do authMiddleware: mais "puro" semanticamente, mas impede logout quando o token já expirou — descartado

---

## 4. Endpoint GET /users/me

**Decisão**: Novo handler `GetMe` no `UserController`. Extrai o `userID` do token via `authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))` e chama `userService.GetByID(ctx, userID)` já existente. Responde com `UserDTO` (mesmo contrato de `GetByID`).

**Rationale**: Reutiliza completamente a lógica existente de busca de usuário. Nenhum novo método de service ou repositório necessário. Registrado ANTES de `GET /users/:userID` nas rotas para que "me" não seja interpretado como `:userID` pelo router.

**Injeção de AuthTokenManager no UserController**: Necessário para extração do userID. O `GroupController` e o `GroupInviteController` já usam este padrão — a mudança é consistente com o projeto.

---

## 5. Configuração de Cookie Seguro por Ambiente

**Decisão**: Adicionar `CookieSecure bool` à struct `AuthConfig` com `envDefault:"true"`. Em `.env` (desenvolvimento), definir `AUTH_COOKIE_SECURE=false`. Passar como parâmetro ao construtor `NewAuthController`.

**Rationale**: Mínimo necessário para satisfazer FR-003 sem over-engineering. `caarlos0/env v11` suporta parsing de bool via `envDefault:"true"`. O parâmetro é um `bool` simples no construtor, sem abstrações adicionais.

**Alternativas consideradas**:
- Detectar automaticamente HTTPS a partir do request: não confiável atrás de proxies reversos — descartado
- Variável de ambiente de ambiente (ex: `APP_ENV=production`): introduz uma abstração extra sem ganho — descartado
