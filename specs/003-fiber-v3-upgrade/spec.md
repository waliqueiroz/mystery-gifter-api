# Feature Specification: Atualização do Go e do Framework Fiber V2 para V3

**Feature Branch**: `003-fiber-v3-upgrade`  
**Created**: 2026-06-21  
**Status**: Draft  

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Atualização do Go para a versão 1.26.4 (Priority: P1)

Como desenvolvedor da API, quero atualizar o toolchain Go para a versão 1.26.4 antes de qualquer outra mudança de dependência, para garantir que as demais atualizações sejam feitas sobre uma base estável e suportada.

**Why this priority**: É pré-requisito das demais stories. Algumas dependências atualizadas nesta feature podem exigir features do compilador ou da stdlib disponíveis apenas em versões recentes do Go. Fazer isso primeiro isola eventuais problemas de toolchain dos problemas de migração de framework.

**Independent Test**: Pode ser testado compilando e executando o projeto com Go 1.26.4 e verificando que todos os testes existentes continuam passando sem alteração.

**Acceptance Scenarios**:

1. **Given** o projeto está declarando Go 1.25.1 em `go.mod`, **When** o toolchain é atualizado para 1.26.4, **Then** o projeto compila sem erros e todos os testes passam.
2. **Given** o toolchain foi atualizado, **When** o comando de build é executado, **Then** o binário é gerado corretamente na versão 1.26.4.

---

### User Story 2 — Todos os endpoints da API continuam funcionando após a atualização (Priority: P2)

Como usuário final da API, quero que todos os endpoints continuem respondendo corretamente após a atualização interna do framework, sem nenhuma mudança no comportamento observável.

**Why this priority**: A continuidade funcional é o critério de aceitação mais crítico. Qualquer regressão introduzida pela migração impacta diretamente os consumidores da API.

**Independent Test**: Pode ser testado executando a suíte de testes atual e verificando que todos os cenários de autenticação, criação de usuários, grupos, convites e sorteios continuam passando sem alterações nos contratos de request/response.

**Acceptance Scenarios**:

1. **Given** a API está em execução com Fiber V3, **When** um cliente envia uma requisição de login com credenciais válidas, **Then** recebe um token JWT e status 200.
2. **Given** a API está em execução com Fiber V3, **When** um cliente envia dados inválidos para qualquer endpoint, **Then** recebe status 422 com a mensagem de erro esperada.
3. **Given** a API está em execução com Fiber V3, **When** um cliente acessa um endpoint protegido sem token JWT, **Then** recebe status 401.
4. **Given** a API está em execução com Fiber V3, **When** um cliente acessa um endpoint inexistente, **Then** recebe status 404.
5. **Given** a API está em execução com Fiber V3, **When** ocorre um erro interno inesperado, **Then** recebe status 500 com a estrutura de erro padrão da API.

---

### User Story 3 — Middleware de autenticação JWT continua operacional (Priority: P3)

Como desenvolvedor da API, quero que o middleware JWT continue validando tokens corretamente com Fiber V3, para que rotas protegidas continuem exigindo autenticação válida.

**Why this priority**: O JWT é a camada de segurança central da API. Uma falha aqui expõe todos os endpoints protegidos.

**Independent Test**: Pode ser testado enviando requisições com tokens válidos, expirados e malformados para rotas protegidas e verificando os status de resposta.

**Acceptance Scenarios**:

1. **Given** a rota está protegida pelo middleware JWT, **When** o cliente envia um token válido no header `Authorization`, **Then** a requisição prossegue normalmente.
2. **Given** a rota está protegida pelo middleware JWT, **When** o cliente envia um token expirado ou inválido, **Then** recebe status 401.
3. **Given** a rota está protegida pelo middleware JWT, **When** o cliente não envia nenhum token, **Then** recebe status 401.

---

### User Story 4 — Middleware de recuperação de pânicos continua ativo (Priority: P4)

Como operador da API, quero que o middleware de recuperação de pânicos (`recover`) continue capturando panics e convertendo-os em respostas 500, para que a API não caia em caso de erros inesperados.

**Why this priority**: Garante estabilidade operacional; sem ele, um panic derrubaria o processo inteiro.

**Independent Test**: Pode ser testado simulando um panic em um handler e verificando que a API retorna 500 sem encerrar o processo.

**Acceptance Scenarios**:

1. **Given** ocorre um panic em um handler, **When** o middleware de recover está configurado, **Then** a API retorna status 500 e continua operando normalmente.

---

### Edge Cases

- O que acontece quando o corpo da requisição está malformado (JSON inválido)? → Deve continuar retornando 422.
- O que acontece com requisições CORS de origens não permitidas? → Deve continuar bloqueando com os headers corretos.
- O que acontece quando o ID de um recurso não é um UUID válido? → Deve retornar 404 (comportamento existente mantido).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O projeto DEVE atualizar a versão declarada do Go em `go.mod` de 1.25.1 para 1.26.4.
- **FR-002**: O sistema DEVE atualizar a dependência `github.com/gofiber/fiber/v2` para `github.com/gofiber/fiber/v3`, mantendo todas as funcionalidades existentes.
- **FR-003**: O sistema DEVE atualizar a dependência `github.com/gofiber/contrib/jwt` para `github.com/gofiber/contrib/v3/jwt`, compatível com Fiber V3.
- **FR-004**: O sistema DEVE adaptar todas as assinaturas de handlers de `*fiber.Ctx` para `fiber.Ctx` (interface), conforme exigido pelo Fiber V3.
- **FR-005**: O sistema DEVE substituir `ctx.BodyParser()` por `ctx.Bind().Body()` em todos os controllers.
- **FR-006**: O sistema DEVE substituir `ctx.QueryParser()` por `ctx.Bind().Query()` em todos os controllers.
- **FR-007**: Nos controllers, o acesso ao token JWT autenticado DEVE ser feito via `jwtware.FromContext(ctx)` (retorna `*jwt.Token` tipado), substituindo `ctx.Locals("user")`. A interface e implementação de `AuthTokenManager.GetAuthUserID(token any)` permanecem inalteradas.
- **FR-008**: O sistema DEVE atualizar a configuração do middleware CORS para usar `[]string` nos campos `AllowOrigins`, `AllowMethods` e `AllowHeaders` (antes eram strings separadas por vírgula).
- **FR-009**: O sistema DEVE atualizar o `CustomErrorHandler` para aceitar `fiber.Ctx` (interface) em vez de `*fiber.Ctx`.
- **FR-010**: O sistema DEVE atualizar todos os imports de `github.com/gofiber/fiber/v2` para `github.com/gofiber/fiber/v3` em todos os arquivos Go (produção e testes).
- **FR-011**: Dependências não relacionadas ao objetivo desta feature DEVEM ser atualizadas somente se falharem a compilar com Go 1.26.4 ou Fiber V3; atualizações oportunistas estão fora do escopo.
- **FR-012**: Todos os testes unitários existentes DEVEM continuar passando sem alteração nos cenários testados.
- **FR-013**: A API DEVE continuar iniciando na porta 8080 com as mesmas configurações de ambiente.

### Inventário de mudanças por arquivo

| Arquivo | Mudanças necessárias |
|---------|----------------------|
| `go.mod` / `go.sum` | Atualizar versão Go para 1.26.4; `fiber/v2` → `fiber/v3`; `contrib/jwt` → `contrib/v3/jwt` |
| `internal/infra/runner.go` | Atualizar imports; adaptar CORS para `[]string`; atualizar recover |
| `internal/infra/entrypoint/middlewares.go` | Atualizar import; `ErrorHandler` de `*fiber.Ctx` → `fiber.Ctx` |
| `internal/infra/entrypoint/routes.go` | Atualizar import |
| `internal/infra/entrypoint/error_handler.go` | Atualizar import; assinaturas `*fiber.Ctx` → `fiber.Ctx` |
| `internal/infra/entrypoint/rest/auth_controller.go` | Atualizar import; `BodyParser` → `Bind().Body()`; `*fiber.Ctx` → `fiber.Ctx` |
| `internal/infra/entrypoint/rest/user_controller.go` | Atualizar import; `BodyParser` → `Bind().Body()`; `QueryParser` → `Bind().Query()`; `*fiber.Ctx` → `fiber.Ctx` |
| `internal/infra/entrypoint/rest/group_controller.go` | Atualizar import; `BodyParser` → `Bind().Body()`; `QueryParser` → `Bind().Query()`; `ctx.Locals("user")` → `jwtware.FromContext(ctx)`; `*fiber.Ctx` → `fiber.Ctx` |
| `internal/infra/entrypoint/rest/group_invite_controller.go` | Atualizar import; `ctx.Locals("user")` → `jwtware.FromContext(ctx)`; `*fiber.Ctx` → `fiber.Ctx` |
| Arquivos `_test.go` nos pacotes `entrypoint` e `rest` | Atualizar imports e assinaturas de `*fiber.Ctx` → `fiber.Ctx` |

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% dos testes unitários existentes passam após a migração sem alteração nos cenários testados.
- **SC-002**: Todos os endpoints da API retornam os mesmos status HTTP e estruturas de resposta que retornavam antes da atualização.
- **SC-003**: A API inicia sem erros de compilação ou de inicialização.
- **SC-004**: Requisições com autenticação JWT válida continuam sendo processadas corretamente em todos os endpoints protegidos.
- **SC-005**: Nenhuma mudança no contrato público da API (paths, métodos HTTP, schemas de request/response, códigos de status).

## Assumptions

- Go 1.26.4 é retrocompatível com o código existente — a atualização do toolchain não exige mudanças no código fonte.
- A versão 1.26.4 do Go satisfaz o requisito mínimo do Fiber V3 (Go 1.25+).
- O pacote `github.com/gofiber/contrib/v3/jwt` é a versão compatível com Fiber V3 e substitui diretamente `github.com/gofiber/contrib/jwt`.
- `jwtware.FromContext(ctx)` do `gofiber/contrib/v3/jwt` substitui `ctx.Locals("user")` para recuperar o token JWT validado; a interface `AuthTokenManager.GetAuthUserID(token any)` permanece inalterada pois `*jwt.Token` satisfaz `any`.
- A API não utiliza os middlewares removidos ou movidos para contrib no Fiber V3 (static, filesystem, session, monitor) — apenas CORS, recover e JWT.
- A migração é puramente interna ao servidor; nenhuma mudança de infraestrutura (Docker, CI/CD, variáveis de ambiente) é necessária.
- O CORS atualmente está configurado com valores padrão permissivos; a adaptação para `[]string` manterá o mesmo comportamento.
- Dependências não relacionadas ao objetivo desta feature (ex: `sqlx`, `squirrel`, `caarlos0/env`) serão atualizadas somente se falharem a compilação com Go 1.26.4 ou Fiber V3; atualizações oportunistas estão fora do escopo.

## Clarifications

### Session 2026-06-21

- Q: Ao atualizar Go para 1.26.4 e Fiber para V3, outras dependências do projeto devem ser atualizadas também? → A: Somente dependências-alvo (Go + Fiber + JWT contrib) mais qualquer outra que falhe a compilar com as novas versões.
