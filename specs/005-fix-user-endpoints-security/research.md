# Research: Correção de Endpoints de Usuário Inseguros

**Branch**: `005-fix-user-endpoints-security` | **Date**: 2026-06-22

## Decisão 1: Remover vs. restringir os endpoints de usuário

**Decisão**: Remover completamente `GET /api/v1/users` e `GET /api/v1/users/:userID`.

**Justificativa**: O único caso de uso legítimo para ver dados de outro usuário (perfil de membro de grupo) já é atendido pela resposta do `GET /api/v1/groups/:groupID`, que inclui todos os membros. O `GET /api/v1/users/me` já cobre o caso do próprio usuário. Manter os endpoints — mesmo restritos — seria superfície de ataque desnecessária.

**Alternativas consideradas**:
- Restringir `GET /users/:id` para permitir apenas `authUserID == userID`: descartado por ser redundante com `/me`.
- Restringir `GET /users` para contexto de grupo (ex: buscar membros dentro de um grupo): descartado por não ter caso de uso no produto atual.

---

## Decisão 2: Como corrigir `GET /api/v1/groups` (user_id + owner_id)

**Decisão**: O handler `GroupController.Search` deve:
1. Extrair `authUserID` do JWT via `c.AuthTokenManager.GetAuthUserID(jwtware.FromContext(ctx))`.
2. Sobrescrever `groupFiltersDTO.UserID = authUserID` antes de mapear para domínio.
3. Se `groupFiltersDTO.OwnerID != ""` e `groupFiltersDTO.OwnerID != authUserID`, retornar `fiber.NewError(fiber.StatusForbidden, "owner_id must match authenticated user")`.

**Justificativa**: A sobrescrita garante que o backend nunca processe um `user_id` arbitrário vindo do cliente. A validação de `owner_id` fecha a segunda brecha sem sacrificar o filtro legítimo de "ver meus grupos onde sou dono".

**Alternativas consideradas**:
- Remover `owner_id` da API: descartado — o filtro tem uso legítimo ("meus grupos onde sou dono").
- Mover validação para a camada de domínio (GroupFilters): descartado — a restrição é de autorização HTTP, não de invariante de domínio.

---

## Decisão 3: Escopo de remoção de código

**Decisão**: Remover em cascata todos os artefatos que dependem exclusivamente dos endpoints extintos:
- `UserController.GetByID`, `UserController.Search`
- `UserService.Search` (interface + implementação)
- `UserRepository.Search` (interface domain + implementação postgres)
- `UserFiltersDTO`, `mapUserFiltersDTOToDomain` (rest/user_dto.go)
- `UserFilters` domain type + `user_filters_builder.go`
- Testes de todas as camadas acima

`UserService.GetByID` e `UserRepository.GetByID` **permanecem** pois são utilizados por `GetMe`, `GroupService.Create` e `GroupService.AddUser`.

**Justificativa**: Princípio VI (YAGNI) — código sem uso presente deve ser removido. Manter código morto aumenta custo de manutenção e confunde futuros colaboradores.

---

## Decisão 4: Atualização da documentação Swagger

**Decisão**: Remover as anotações `swagger:operation` dos dois endpoints extintos em `routes.go`. Atualizar a anotação de `GET /api/v1/groups` para:
- Remover o parâmetro `user_id` da lista de query params (agora ignorado e sobrescrito).
- Atualizar a descrição de `owner_id` para indicar que deve ser o ID do usuário autenticado.

**Justificativa**: Documentação inconsistente com a implementação cria confusão para consumidores futuros e viola o gate de qualidade da constituição (`make generate-docs` deve passar e refletir a realidade).
