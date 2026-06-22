# Quickstart de Desenvolvimento: Autenticação Dual-Channel com Cookie HttpOnly

## Pré-requisitos

- Go 1.26.4+
- Docker (PostgreSQL local)
- `.env` baseado em `.env.example`

## Configuração Local

Adicionar ao `.env`:

```env
AUTH_COOKIE_SECURE=false
```

O valor `false` é necessário para desenvolvimento local (HTTP). Em produção, a variável deve ser `true` (valor padrão).

## Rodar a Aplicação

```bash
make run
```

## Verificar a Feature Localmente

### 1. Login via cookie

```bash
curl -c cookies.txt -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "mypassword"}'
```

O arquivo `cookies.txt` armazenará o cookie `access_token`. Verificar que o `Set-Cookie` está presente na resposta com os atributos `HttpOnly` e `SameSite=Lax`.

### 2. Acessar rota protegida via cookie

```bash
curl -b cookies.txt http://localhost:8080/api/v1/users/me
```

### 3. Acessar rota protegida via cabeçalho (retrocompatibilidade mobile)

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/users/me
```

### 4. Logout

```bash
curl -c cookies.txt -b cookies.txt -X POST http://localhost:8080/api/v1/logout
```

Verificar que o cookie é removido (header `Set-Cookie: access_token=; MaxAge=-1`).

### 5. Acesso após logout (deve retornar 400)

```bash
curl -b cookies.txt http://localhost:8080/api/v1/users/me
```

## Rodar os Testes

```bash
make test
```

Todos os testes devem passar. Verificar cobertura dos novos handlers:
- `Test_AuthController_Login` — novo cenário: cookie definido na resposta
- `Test_AuthController_Logout` — novo
- `Test_UserController_GetMe` — novo
- `Test_NewAuthMiddleware` — novos cenários: cookie válido, cookie inválido, ambos presentes

## Gerar Documentação

```bash
make generate-docs
make serve-docs  # acessa http://localhost:8081
```

Verificar os novos endpoints `POST /logout` e `GET /users/me` na UI do Swagger.
