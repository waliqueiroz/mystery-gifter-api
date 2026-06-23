# Implementation Plan: Correção de Endpoints de Usuário Inseguros

**Branch**: `005-fix-user-endpoints-security` | **Date**: 2026-06-22 | **Spec**: [spec.md](spec.md)

## Summary

Remoção dos endpoints `GET /api/v1/users` e `GET /api/v1/users/:userID` e de toda a infraestrutura associada (handler, DTOs, service method, repository method, mocks, builders, testes). Correção do handler `GET /api/v1/groups` para sobrescrever o parâmetro `user_id` com o `authUserID` do token JWT e rejeitar requisições onde `owner_id` não corresponde ao usuário autenticado.

## Technical Context

**Language/Version**: Go 1.26.4
**Primary Dependencies**: Fiber v3, gofiber/contrib/v3/jwt, golang-jwt v5, go-playground/validator, go.uber.org/mock/mockgen, go-swagger
**Storage**: PostgreSQL via sqlx + squirrel
**Testing**: testify/assert + go.uber.org/mock/mockgen
**Target Platform**: Linux server
**Project Type**: web-service (REST API)
**Performance Goals**: p95 < 200ms para CRUD padrão (sem impacto desta mudança)
**Constraints**: Clean Architecture estrita — camadas domain → application → infra; testes unitários obrigatórios por ciclo de implementação

## Constitution Check

| Princípio | Status | Observação |
|---|---|---|
| I. Clean Architecture | ✅ Passa | Remoção em cascata em todas as camadas, respeitando a direção de dependência |
| II. Disciplina de Testes | ✅ Passa | Testes dos handlers/service/repo removidos serão excluídos; testes de `GroupController.Search` serão atualizados |
| III. Validação no Domínio | ✅ Passa | Nenhuma lógica de domínio nova; apenas remoção e extração de parâmetro de auth |
| IV. Contrato de API Consistente | ✅ Passa | Handler `GroupController.Search` atualizado com extração de `authUserID`; Swagger atualizado via `make generate-docs` |
| V. Abstração de Infraestrutura | ✅ Passa | `UserRepository.Search` removido da interface → mock regenerado via `go generate` |
| VI. Simplicidade & YAGNI | ✅ Passa | Remoção líquida: sem abstrações novas, apenas eliminação de código não mais justificado |
| VII. Performance & Observabilidade | ✅ Passa | Sem impacto; paginação e ordenação mantidas no `GET /groups` |
| VIII. Idioma dos Artefatos | ✅ Passa | Artefatos do speckit em pt-BR; código, testes e mensagens de erro em inglês |

## Project Structure

### Documentation (this feature)

```text
specs/005-fix-user-endpoints-security/
├── plan.md              # Este arquivo
├── research.md          # Fase 0
├── contracts/
│   └── api.md           # Contrato de endpoints atualizado
└── tasks.md             # Gerado por /speckit.tasks
```

### Source Code (repository root)

```text
internal/
├── domain/
│   ├── user.go                           # Remover UserRepository.Search da interface
│   ├── mock_domain/
│   │   └── user_repository.go            # Regenerar: go generate ./internal/domain/...
│   └── build_domain/
│       └── user_filters_builder.go       # Excluir (UserFilters não mais utilizado)
├── application/
│   ├── user_service.go                   # Remover Search da interface + implementação
│   ├── user_service_test.go              # Remover testes do método Search
│   └── mock_application/
│       └── user_service.go               # Regenerar: go generate ./internal/application/...
└── infra/
    ├── entrypoint/
    │   ├── routes.go                     # Remover rotas + Swagger de GET /users e GET /users/:userID;
    │   │                                 # atualizar Swagger de GET /groups (remover user_id, documentar owner_id restrito)
    │   └── rest/
    │       ├── user_controller.go        # Remover métodos GetByID e Search
    │       ├── user_controller_test.go   # Remover testes de GetByID e Search
    │       ├── user_dto.go               # Remover UserFiltersDTO, mapUserFiltersDTOToDomain
    │       ├── group_controller.go       # Corrigir Search: extrair authUserID, sobrescrever UserID, validar OwnerID
    │       └── group_controller_test.go  # Atualizar testes de Search
    └── outgoing/
        └── postgres/
            ├── user_repository.go        # Remover implementação de Search
            └── user_repository_test.go   # Remover testes de Search
```
