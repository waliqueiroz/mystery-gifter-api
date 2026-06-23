# API Contract: Endpoints Afetados

**Feature**: `005-fix-user-endpoints-security` | **Date**: 2026-06-22

## Endpoints Removidos

### ~~GET /api/v1/users~~
**Status**: REMOVIDO — retorna 404 após deploy desta feature.

### ~~GET /api/v1/users/:userID~~
**Status**: REMOVIDO — retorna 404 após deploy desta feature.

---

## Endpoints Inalterados

### GET /api/v1/users/me
**Status**: Sem alteração. Continua retornando os dados do usuário autenticado.

---

## Endpoints Modificados

### GET /api/v1/groups

**Antes**: Aceitava `user_id` e `owner_id` como query params sem validação de autorização.

**Depois**:

| Parâmetro | Comportamento |
|-----------|--------------|
| `user_id` | Ignorado. O backend sempre usa o `authUserID` extraído do JWT. |
| `owner_id` | Aceito, mas deve ser igual ao `authUserID`. Caso contrário: `403 Forbidden`. |
| `name`, `status[]`, `limit`, `offset`, `sort_by`, `sort_direction` | Sem alteração. |

**Novo comportamento de erro**:

```
403 Forbidden
{
  "code": "forbidden",
  "message": "owner_id must match authenticated user"
}
```

**Exemplo de requisição válida**:
```
GET /api/v1/groups?status=OPEN&status=MATCHED&sort_direction=DESC&limit=15&offset=0
Authorization: Bearer <token>
```

**Exemplo de requisição que retorna 403**:
```
GET /api/v1/groups?owner_id={idDeOutraPessoa}
Authorization: Bearer <token>
```
