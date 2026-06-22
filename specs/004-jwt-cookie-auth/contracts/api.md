# Contratos de API: Autenticação Dual-Channel com Cookie HttpOnly

Base URL: `/api/v1`

---

## POST /login (atualizado)

Autentica o usuário e define o cookie de sessão **além** de retornar o token no corpo.

**Autenticação**: Não requerida  
**Alteração em relação ao comportamento atual**: Adição do Set-Cookie na resposta

### Request

```
POST /api/v1/login
Content-Type: application/json
```

```json
{
  "email": "user@example.com",
  "password": "mypassword123"
}
```

### Response (200 OK) — sem alteração no corpo

```json
{
  "user": {
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "name": "João",
    "surname": "Silva",
    "email": "user@example.com",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 1735689600
}
```

**Novo Set-Cookie no header de resposta:**

```
Set-Cookie: access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...; 
            Path=/; 
            Expires=Wed, 01 Jan 2025 00:00:00 GMT; 
            HttpOnly; 
            Secure; 
            SameSite=Lax
```

_(`Secure` omitido quando `AUTH_COOKIE_SECURE=false`)_

### Respostas de Erro (sem alteração)

| Status | Código         | Quando                              |
|--------|----------------|-------------------------------------|
| 400    | bad_request    | Credenciais inválidas ou campo obrigatório ausente |
| 401    | unauthorized   | Email ou senha incorretos           |
| 422    | unprocessable  | Payload malformado                  |

---

## POST /logout (novo)

Remove o cookie de sessão do navegador. **Não requer autenticação** — opera apenas como limpeza de cookie no cliente.

**Autenticação**: Não requerida (registrado antes do `authMiddleware`)

### Request

```
POST /api/v1/logout
```

Sem corpo.

### Response (204 No Content)

Sem corpo.

**Set-Cookie no header de resposta (remoção):**

```
Set-Cookie: access_token=; 
            Path=/; 
            Expires=Thu, 01 Jan 1970 00:00:00 GMT; 
            MaxAge=-1; 
            HttpOnly; 
            SameSite=Lax
```

---

## GET /users/me (novo)

Retorna os dados do usuário atualmente autenticado. Mesmo contrato de `GET /users/:userID`.

**Autenticação**: Requerida (cookie ou cabeçalho `Authorization: Bearer`)  
**Registrado**: ANTES de `GET /users/:userID` para evitar conflito de roteamento

### Request

```
GET /api/v1/users/me
Cookie: access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Ou via cabeçalho:

```
GET /api/v1/users/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Response (200 OK)

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "name": "João",
  "surname": "Silva",
  "email": "user@example.com",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### Respostas de Erro

| Status | Código       | Quando                              |
|--------|--------------|-------------------------------------|
| 400    | bad_request  | Token ausente (nem cookie nem header) |
| 401    | unauthorized | Token inválido ou expirado          |
| 404    | not_found    | Usuário não encontrado (token válido, mas usuário deletado) |

---

## Middleware de Autenticação (atualizado)

Todas as rotas protegidas usam o middleware com a seguinte lógica de extração:

```
1. Tenta ler cookie "access_token"
   → Encontrado: valida o JWT do cookie
     → Válido: autoriza (não tenta cabeçalho)
     → Inválido: retorna 401 (não faz fallback para cabeçalho)
2. Cookie ausente: tenta ler "Authorization: Bearer <token>"
   → Encontrado e válido: autoriza
   → Encontrado e inválido: retorna 401
   → Ausente: retorna 400 (missing or malformed JWT)
```

**Implementação**: `extractors.Chain(extractors.FromCookie("access_token"), extractors.FromAuthHeader("Bearer"))` via `jwtware.Config.Extractor`.
