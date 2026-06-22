# Modelo de Dados: Autenticação Dual-Channel com Cookie HttpOnly

## Entidades Existentes (sem alteração de schema)

Esta feature **não requer migrações de banco de dados**. O cookie de autenticação é um mecanismo puramente HTTP — não há persistência de sessão no servidor.

### AuthSession (domínio — sem alterações)

| Campo       | Tipo     | Descrição                                                  |
|-------------|----------|------------------------------------------------------------|
| User        | User     | Dados completos do usuário autenticado                     |
| AccessToken | string   | Token JWT assinado                                         |
| TokenType   | string   | Tipo do token (ex: "Bearer")                               |
| ExpiresIn   | int64    | Unix timestamp de expiração — reutilizado como Expires do cookie |

### User (domínio — sem alterações)

Campos públicos expostos via `UserDTO`: `ID`, `Name`, `Surname`, `Email`, `CreatedAt`, `UpdatedAt`. O campo `Password` (hash) nunca é incluído nas respostas.

---

## Novos Elementos de Infraestrutura HTTP

### Cookie de Autenticação

| Atributo  | Valor                                          |
|-----------|------------------------------------------------|
| Name      | `access_token`                                 |
| Value     | JWT assinado (mesmo valor de `AuthSessionDTO.access_token`) |
| Expires   | `time.Unix(authSession.ExpiresIn, 0)`          |
| HTTPOnly  | `true` (sempre)                                |
| Secure    | Configurável via `AUTH_COOKIE_SECURE` (default: `true`) |
| SameSite  | `Lax`                                          |
| Path      | `/` (padrão)                                   |

**Remoção do cookie (logout)**: Definir `Expires` no passado (ex: `time.Unix(0, 0)`) e `MaxAge: -1`.

---

## Variável de Ambiente Nova

| Variável             | Tipo   | Default | Descrição                                     |
|----------------------|--------|---------|-----------------------------------------------|
| `AUTH_COOKIE_SECURE` | bool   | `true`  | Define flag Secure do cookie; `false` em dev  |

Esta variável deve ser adicionada ao `AuthConfig` e ao `.env.example`.
