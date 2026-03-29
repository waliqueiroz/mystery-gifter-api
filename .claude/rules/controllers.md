---
paths:
  - "internal/infra/entrypoint/rest/**/*.go"
---

# Convenções de Controllers e Mappers

## Controllers

- Pacote `rest`, struct `XController`, construtor `NewXController`
- Handler: `func (c *XController) Method(ctx *fiber.Ctx) error`
- Fluxo obrigatório: `BodyParser` → `dto.Validate()` → `mapXToDomain` → service → `mapXFromDomain` → `ctx.JSON`
- `BodyParser` falha → `fiber.NewError(fiber.StatusUnprocessableEntity)`; demais erros → retornar diretamente
- Criação bem-sucedida → `ctx.Status(fiber.StatusCreated).JSON(...)`; outras operações → `ctx.JSON(...)`
- Parâmetros de rota via `ctx.Params("paramName")`, nunca `ctx.Query` para IDs de rota
- Auth user ID: `authUserID, err := c.AuthTokenManager.GetAuthUserID(ctx.Locals("user"))`
- Nenhuma lógica de negócio nos controllers

## Mappers

- Funções privadas, no mesmo arquivo do tipo que servem
- `mapXToDomain` para DTO/persistence → domínio; `mapXFromDomain` para domínio → DTO/persistence
- Slices no plural: `mapUsersFromDomain`, `mapUsersToDomain`
- Assinatura sempre retorna `(result, error)`
- Dependências externas (IdentityGenerator, PasswordManager) como primeiros argumentos
- Chamar `Validate()` na entrada (DTO) e na saída (objeto criado), retornar `nil, err` se falhar
- Slices: pré-alocar com `make([]T, 0, len(src))`, propagar erros imediatamente
