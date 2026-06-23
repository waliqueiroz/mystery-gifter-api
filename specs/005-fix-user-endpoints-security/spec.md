# Feature Specification: Correção de Endpoints de Usuário Inseguros

**Feature Branch**: `005-fix-user-endpoints-security`
**Created**: 2026-06-22
**Status**: Done
**Input**: User description: "Corrigir brechas de segurança nos endpoints de usuário: remover GET /users e GET /users/:id (substituídos por /me), e forçar filtragem de grupos pelo usuário autenticado em GET /groups"

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Usuário não consegue acessar dados pessoais de outros usuários (Priority: P1)

Um usuário autenticado tenta acessar informações de outro usuário diretamente pela URL. Com a remoção dos endpoints `GET /users` e `GET /users/:id`, o sistema simplesmente não oferece mais essa superfície — qualquer tentativa resulta em erro 404.

**Why this priority**: É a brecha mais grave: permite que qualquer usuário autenticado enumere todos os usuários do sistema e acesse e-mail, nome e sobrenome de qualquer pessoa.

**Independent Test**: Pode ser testado isoladamente realizando requisições para `GET /api/v1/users` e `GET /api/v1/users/:id` com um token válido e verificando que ambas retornam 404.

**Acceptance Scenarios**:

1. **Given** um usuário autenticado, **When** ele realiza `GET /api/v1/users`, **Then** o sistema retorna 404.
2. **Given** um usuário autenticado, **When** ele realiza `GET /api/v1/users/{qualquerID}`, **Then** o sistema retorna 404.
3. **Given** um usuário autenticado, **When** ele realiza `GET /api/v1/users/me`, **Then** o sistema retorna os dados do próprio usuário normalmente.

---

### User Story 2 — Usuário só visualiza grupos aos quais pertence (Priority: P1)

Um usuário autenticado solicita a listagem de grupos. O sistema deve retornar exclusivamente os grupos onde o usuário autenticado é membro, independentemente de qualquer parâmetro enviado na requisição.

**Why this priority**: Sem essa correção, um atacante pode passar o ID de outra pessoa como `user_id` na query string e visualizar todos os grupos daquela pessoa.

**Independent Test**: Pode ser testado isoladamente chamando `GET /api/v1/groups` com e sem o parâmetro `user_id`, verificando que o resultado sempre reflete apenas os grupos do usuário autenticado.

**Acceptance Scenarios**:

1. **Given** um usuário autenticado que pertence a 3 grupos, **When** ele realiza `GET /api/v1/groups`, **Then** o sistema retorna apenas esses 3 grupos.
2. **Given** um usuário autenticado, **When** ele realiza `GET /api/v1/groups?user_id={IDDeOutroUsuário}`, **Then** o sistema ignora o parâmetro e retorna apenas os grupos do usuário autenticado.
3. **Given** um usuário autenticado sem grupos, **When** ele realiza `GET /api/v1/groups`, **Then** o sistema retorna lista vazia.
4. **Given** uma requisição sem token de autenticação, **When** ela acessa `GET /api/v1/groups`, **Then** o sistema retorna 401.

---

### Edge Cases

- O que acontece se o usuário enviar `user_id` de si mesmo no parâmetro? Deve funcionar normalmente (o resultado é o mesmo).
- O que acontece com os demais filtros de `GET /api/v1/groups` (`name`, `status[]`, `sort_by`, `sort_direction`, `limit`, `offset`)? Devem continuar funcionando normalmente.
- O que acontece com `owner_id` como filtro? Se enviado, deve ser validado contra o `authUserID` — retorna 403 se não corresponder.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O sistema DEVE remover o endpoint `GET /api/v1/users`, tornando-o inacessível (404) para qualquer chamada.
- **FR-002**: O sistema DEVE remover o endpoint `GET /api/v1/users/:userID`, tornando-o inacessível (404) para qualquer chamada.
- **FR-003**: O endpoint `GET /api/v1/users/me` DEVE continuar funcionando e retornando os dados do usuário autenticado.
- **FR-004**: O endpoint `GET /api/v1/groups` DEVE ignorar o parâmetro `user_id` enviado pelo cliente e sempre usar o identificador do usuário autenticado no token para filtrar grupos.
- **FR-005**: Os filtros de `GET /api/v1/groups` (`name`, `status[]`, paginação, ordenação) DEVEM continuar funcionando normalmente.
- **FR-008**: Se o parâmetro `owner_id` for enviado em `GET /api/v1/groups`, o sistema DEVE validar que ele corresponde ao `authUserID`; caso contrário, retornar 403.
- **FR-006**: O sistema DEVE remover toda lógica de serviço e repositório que não é mais utilizada após a remoção dos endpoints.
- **FR-007**: A documentação Swagger/OpenAPI DEVE ser atualizada para remover completamente as definições dos endpoints extintos (`GET /users` e `GET /users/:userID`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Chamadas para `GET /api/v1/users` e `GET /api/v1/users/:id` retornam 404 em 100% dos casos, independentemente do token apresentado.
- **SC-002**: Chamadas para `GET /api/v1/groups` com `user_id` de outro usuário retornam exclusivamente os grupos do usuário autenticado, nunca os do usuário informado no parâmetro.
- **SC-003**: Nenhuma regressão nos endpoints restantes — todos os testes existentes continuam passando.
- **SC-004**: O endpoint `GET /api/v1/users/me` continua retornando dados corretos do usuário autenticado em 100% das chamadas válidas.

## Clarifications

### Session 2026-06-22

- Q: A documentação Swagger deve ser atualizada para refletir a remoção dos endpoints? → A: Sim, remover completamente as definições dos endpoints extintos da documentação Swagger/OpenAPI.
- Q: O filtro `owner_id` em `GET /groups` deve ser restringido? → A: Sim — se enviado, deve obrigatoriamente ser o `authUserID`; caso contrário, retornar 403.
- Q: Qual a ordem de deploy entre frontend e backend? → A: Frontend primeiro — a remoção das chamadas de API deve estar em produção antes dos endpoints serem removidos do backend.

## Assumptions

- O frontend será atualizado em paralelo para não depender mais dos endpoints removidos.
- Não há outro consumidor externo (mobile, integração terceira) que dependa de `GET /users` ou `GET /users/:id`.
- O parâmetro `owner_id` em `GET /groups` não representa risco equivalente, pois retorna grupos filtrados por dono sem expor dados pessoais.
- Os métodos de serviço e repositório correspondentes aos endpoints removidos podem ser excluídos sem impacto em outros fluxos.
