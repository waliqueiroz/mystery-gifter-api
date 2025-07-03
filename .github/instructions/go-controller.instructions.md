---
description: Padrões e convenções para implementação de controllers REST no projeto mystery-gifter-api
applyTo: internal/infra/entrypoint/rest/*controller.go
---
# Convenções para Controllers REST (`*controller.go`)

Este guia define os padrões obrigatórios para implementação de controllers REST no projeto `mystery-gifter-api`. Siga estas diretrizes para garantir consistência, clareza e alinhamento com a arquitetura do projeto.

## Estrutura Geral

- **Pacote:** Todos os controllers devem estar no pacote `rest`.
- **Imports:** Importe sempre o Fiber, os serviços da camada de aplicação e as interfaces de domínio relevantes.
- **Nome do Controller:** Use o padrão `XController` (ex: `UserController`, `GroupController`, `AuthController`).

## Injeção de Dependências

- **Construtor:** Implemente sempre uma função `NewXController` para injeção explícita das dependências (serviços, managers, etc).
- **Campos:** Armazene as dependências como campos privados do struct do controller.

```go
type UserController struct {
	userService       application.UserService
	identityGenerator domain.IdentityGenerator
	passwordManager   domain.PasswordManager
}
```

## Métodos de Handler

- **Assinatura:** Todos os handlers devem ter a assinatura `func (c *XController) Method(ctx *fiber.Ctx) error`.
- **Nomenclatura:** Use nomes claros e descritivos, geralmente ações HTTP (`Create`, `GetByID`, `AddUser`, `RemoveUser`, `Login`).

## Processamento de Requisições

1. **DTOs:** 
   - Sempre utilize DTOs para entrada e saída de dados.
   - Faça o parsing do corpo da requisição para o DTO correspondente usando `ctx.BodyParser(&dto)`.
   - Valide o DTO com um método `Validate()` quando aplicável.

2. **Validação de Erros:**
   - Se o parsing do corpo falhar, retorne `fiber.NewError(fiber.StatusUnprocessableEntity)`.
   - Propague erros de validação e de serviço diretamente (`return err`).

3. **Autenticação:**
   - Quando necessário, recupere o ID do usuário autenticado via `AuthTokenManager.GetAuthUserID(ctx.Locals("user"))`.

4. **Chamada ao Serviço:**
   - Delegue a lógica de negócio para o serviço da camada de aplicação, passando sempre o contexto (`ctx.Context()`).

5. **Mapeamento de Resposta:**
   - Converta entidades de domínio para DTOs de resposta usando funções de mapeamento (`mapXFromDomain`).

6. **Resposta HTTP:**
   - Use `ctx.Status(fiber.StatusCreated).JSON(dto)` para respostas de criação.
   - Use `ctx.JSON(dto)` para respostas de consulta.

## Exemplo de Handler

```go
func (c *UserController) Create(ctx *fiber.Ctx) error {
	var createUserDTO CreateUserDTO

	if err := ctx.BodyParser(&createUserDTO); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	user, err := mapCreateUserDTOToDomain(c.identityGenerator, c.passwordManager, createUserDTO)
	if err != nil {
		return err
	}

	err = c.userService.Create(ctx.Context(), *user)
	if err != nil {
		return err
	}

	userDTO, err := mapUserFromDomain(*user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(userDTO)
}
```

## Boas Práticas

- **Não implemente lógica de negócio nos controllers.** Use sempre os serviços da camada de aplicação.
- **Propague erros de forma clara e direta.**
- **Mantenha handlers pequenos e focados em orquestrar parsing, validação, chamada de serviço e resposta.**
- **Utilize funções auxiliares para mapeamento entre DTOs e entidades de domínio.**
- **Evite duplicação de código entre handlers semelhantes.**

---

> Consulte esta regra ao criar ou modificar qualquer controller REST para garantir alinhamento com os padrões do projeto.
