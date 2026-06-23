# Tasks: Correção de Endpoints de Usuário Inseguros

**Input**: Design documents from `/specs/005-fix-user-endpoints-security/`
**Branch**: `005-fix-user-endpoints-security`

## Format: `[ID] [P?] [Story] Descrição`

- **[P]**: Pode rodar em paralelo (arquivos diferentes, sem dependências)
- **[Story]**: A qual user story a tarefa pertence (US1, US2)

---

## Phase 1: Foundational — Remoção das interfaces (pré-requisito bloqueante)

**Propósito**: Remover `Search` das interfaces de domínio e serviço antes de qualquer implementação. Isso garante que o compilador do Go aponte todos os lugares que precisam ser atualizados.

**⚠️ CRÍTICO**: Nenhuma tarefa das User Stories pode começar antes desta fase.

- [x] T001 Remover método `Search` da interface `UserRepository` em `internal/domain/user.go`
- [x] T002 Remover método `Search` da interface `UserService` e sua implementação em `internal/application/user_service.go`
- [x] T003 Executar `go generate ./internal/domain/... && go generate ./internal/application/...` para regenerar os mocks `internal/domain/mock_domain/user_repository.go` e `internal/application/mock_application/user_service.go`
- [x] T004 [P] Excluir o arquivo `internal/domain/build_domain/user_filters_builder.go` (não mais utilizado após remoção de `UserFilters`)

**Checkpoint**: `go build ./...` deve falhar apontando para os arquivos que ainda referenciam `Search` — isso é esperado e indica onde as próximas tarefas devem atuar.

---

## Phase 2: User Story 1 — Remoção dos endpoints de usuário (Priority: P1) 🎯 MVP

**Goal**: `GET /api/v1/users` e `GET /api/v1/users/:userID` deixam de existir. Qualquer chamada retorna 404.

**Independent Test**: Fazer `GET /api/v1/users` e `GET /api/v1/users/{qualquer-uuid}` com um token JWT válido e verificar que ambos retornam 404.

- [x] T005 [P] [US1] Remover implementação do método `Search` de `internal/infra/outgoing/postgres/user_repository.go`
- [x] T006 [P] [US1] Remover `UserFiltersDTO`, `mapUserFiltersDTOToDomain` e o tipo `UserSearchResultDTO` (swagger model) de `internal/infra/entrypoint/rest/user_dto.go`
- [x] T007 [US1] Remover os métodos `GetByID` e `Search` de `internal/infra/entrypoint/rest/user_controller.go` (T005 e T006 devem estar concluídos)
- [x] T008 [US1] Remover as rotas `GET /users` e `GET /users/:userID` e seus blocos de anotação Swagger correspondentes de `internal/infra/entrypoint/routes.go` (T007 deve estar concluído)
- [x] T009 [P] [US1] Remover testes dos métodos `GetByID` e `Search` de `internal/infra/entrypoint/rest/user_controller_test.go`
- [x] T010 [P] [US1] Remover testes do método `Search` de `internal/application/user_service_test.go`
- [x] T011 [P] [US1] Remover testes do método `Search` de `internal/infra/outgoing/postgres/user_repository_test.go`
- [x] T012 [US1] Executar `make test` e corrigir qualquer falha antes de prosseguir

**Checkpoint**: `GET /api/v1/users` e `GET /api/v1/users/:userID` retornam 404. `GET /api/v1/users/me` continua funcionando. `make test` passa.

---

## Phase 3: User Story 2 — Correção do filtro de grupos (Priority: P1)

**Goal**: `GET /api/v1/groups` sempre filtra pelo `authUserID` do token JWT. Se `owner_id` for enviado e não corresponder ao usuário autenticado, retorna 403.

**Independent Test**: Fazer `GET /api/v1/groups?user_id={idDeOutroUsuario}` e verificar que retorna apenas os grupos do usuário autenticado. Fazer `GET /api/v1/groups?owner_id={idDeOutroUsuario}` e verificar que retorna 403.

- [x] T013 [US2] Atualizar `GroupController.Search` em `internal/infra/entrypoint/rest/group_controller.go`: extrair `authUserID` via `c.AuthTokenManager.GetAuthUserID(jwtware.FromContext(ctx))`, sobrescrever `groupFiltersDTO.UserID = authUserID`, e retornar `fiber.NewError(fiber.StatusForbidden, "owner_id must match authenticated user")` se `groupFiltersDTO.OwnerID != "" && groupFiltersDTO.OwnerID != authUserID`
- [x] T014 [US2] Atualizar anotação Swagger de `GET /api/v1/groups` em `internal/infra/entrypoint/routes.go`: remover parâmetro `user_id`; atualizar descrição de `owner_id` para indicar que deve ser o ID do usuário autenticado, caso contrário 403
- [x] T015 [US2] Atualizar testes de `GroupController.Search` em `internal/infra/entrypoint/rest/group_controller_test.go`: adicionar casos que verificam que `user_id` da query é ignorado, que `authUserID` é sempre usado, e que `owner_id` diferente do `authUserID` retorna 403
- [x] T016 [US2] Executar `make test` e corrigir qualquer falha antes de prosseguir

**Checkpoint**: `make test` passa. O handler de grupos usa o token como fonte de verdade para filtragem.

---

## Phase 4: Polish & Validação Final

**Propósito**: Garantir que todos os gates de qualidade da constituição passam.

- [x] T017 Executar `make generate-docs` e verificar que a spec Swagger gerada não contém mais os endpoints `GET /users` e `GET /users/:userID` e que `GET /groups` não lista `user_id` como parâmetro aceito
- [x] T018 Executar `make build` e confirmar compilação limpa sem warnings

---

## Dependencies & Execution Order

### Dependências entre phases

- **Phase 1 (Foundational)**: Sem dependências — começar imediatamente
- **Phase 2 (US1)**: Depende do Phase 1 completo
- **Phase 3 (US2)**: Pode começar em paralelo com Phase 2 (arquivos diferentes: group_controller vs user_controller)
- **Phase 4 (Polish)**: Depende de Phase 2 e Phase 3 completos

### Dentro de cada User Story

- T005 e T006 podem rodar em paralelo (arquivos diferentes)
- T007 depende de T005 e T006 (remove handler que usa UserFiltersDTO e user_repository.Search)
- T008 depende de T007 (remove rota que referencia o handler)
- T009, T010, T011 podem rodar em paralelo entre si e com T005/T006

### Oportunidades de paralelismo

```bash
# Phase 2 — paralelo inicial:
T005: Remover Search de user_repository.go
T006: Remover UserFiltersDTO de user_dto.go
T009: Remover testes de user_controller_test.go
T010: Remover testes de user_service_test.go
T011: Remover testes de user_repository_test.go

# Phase 3 — pode rodar junto com Phase 2:
T013: Corrigir GroupController.Search
T015: Atualizar testes de group_controller_test.go
```

---

## Implementation Strategy

### MVP (US1 apenas)

1. Completar Phase 1 (Foundational)
2. Completar Phase 2 (US1)
3. **Validar**: `GET /users` e `GET /users/:id` retornam 404; `/me` funciona; `make test` passa
4. Commitar

### Entrega Incremental

1. Phase 1 → Phase 2 → Validar US1 → Commitar
2. Phase 3 → Validar US2 → Commitar
3. Phase 4 → Commitar e abrir PR

---

## Notes

- [P] = arquivos diferentes, sem dependências entre si
- Após T003 (`go generate`), o compilador vai apontar todos os locais que ainda referenciam `Search` — use como guia
- Manter `UserService.GetByID` e `UserRepository.GetByID`: ainda usados por `GetMe`, `GroupService.Create` e `GroupService.AddUser`
- `make test` deve ser executado ao final de cada User Story antes de avançar
