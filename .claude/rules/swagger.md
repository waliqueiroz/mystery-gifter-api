---
paths:
  - "docs/**/*.go"
  - "internal/infra/entrypoint/routes.go"
  - "internal/infra/entrypoint/rest/**/*.go"
---

# Convenções de Documentação Swagger (go-swagger)

## Arquivos

- `docs/swagger.go` — metadados globais da API (`// swagger:meta`)
- `docs/specs/swagger.yaml` — especificação gerada automaticamente via `make generate-docs`
- `internal/infra/entrypoint/routes.go` — anotações de endpoints

## DTOs

- Adicionar `// swagger:model NomeDTO` no struct
- Documentar cada campo com comentários: `// required: true/false`, `// example: value`, `// enum: A,B,C`, `// minLength: N`
- DTOs de filtros/paginação: usar `query` tag junto com `json`

## Endpoints (em routes.go)

```go
// swagger:operation METHOD /api/v1/path OperationName
// ---
// tags:
// - tag_name
// security:
// - Bearer: []
// parameters:
// responses:
//   '200':
//     schema:
//       "$ref": '#/definitions/ResponseDTO'
```

- Tags disponíveis: `auth`, `users`, `groups`
- Códigos: 200 (GET/PUT/DELETE), 201 (POST criação), 400, 401, 403, 404, 409, 422

## CRÍTICO: Query Parameters

**NUNCA** usar `schema: "$ref"` para query parameters — go-swagger + OpenAPI 2.0 não renderiza corretamente.

Sempre definir cada query parameter individualmente:
```go
// - name: fieldName
//   in: query
//   description: descrição
//   required: false
//   type: string
```

Body parameters podem (e devem) usar `schema: "$ref"`.
