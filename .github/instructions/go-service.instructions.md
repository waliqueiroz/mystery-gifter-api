---
description: Convenções para implementação de serviços (Service Layer) no projeto mystery-gifter-api
applyTo: internal/application/**/*_service.go
---
# Convenções para Implementação de Serviços (Service Layer)

Este documento define o padrão para criação e manutenção de serviços na camada de aplicação do projeto `mystery-gifter-api`. Siga estas diretrizes para garantir consistência, testabilidade e clareza.

1. **Interface e Implementação**
   - Defina uma interface para cada serviço, nomeada como `<Entity>Service` (ex: `UserService`, `AuthService`, `GroupService`).
   - Implemente a interface em uma struct não exportada (ex: `type userService struct`).
   - Use o padrão Factory: forneça uma função construtora (ex: `NewUserService(...) UserService`) que retorna a interface.

```go
type UserService interface {
    Create(ctx context.Context, user domain.User) error
    GetByID(ctx context.Context, userID string) (*domain.User, error)
}

type userService struct {
    userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) UserService {
    return &userService{userRepository: userRepository}
}
```

2. **Injeção de Dependências**
   - Receba dependências via construtor (ex: repositórios, geradores, managers).
   - Nunca acople dependências diretamente dentro da implementação.

3. **Contexto**
   - Todos os métodos públicos devem receber `context.Context` como primeiro parâmetro.
   - Propague o contexto para chamadas de repositórios e dependências.

4. **Validação**
   - Valide entidades de entrada usando o método `Validate()` antes de prosseguir.
   - Retorne erro imediatamente se a validação falhar.

```go
if err := user.Validate(); err != nil {
    return err
}
```

5. **Tratamento de Erros**
   - Padronize mensagens de erro usando funções de erro do domínio (ex: `domain.NewUnauthorizedError`).
   - Nunca exponha detalhes sensíveis em mensagens de erro.

6. **Retorno**
   - Retorne ponteiros para structs quando apropriado (ex: `(*domain.User, error)`).
   - Para métodos de criação, retorne a entidade criada ou erro.

7. **Testabilidade**
   - Use diretiva `//go:generate` para gerar mocks das interfaces de serviço.
   - Implemente testes unitários para cada serviço.

```go
//go:generate go run go.uber.org/mock/mockgen -destination mock_application/user_service.go . UserService
```

8. **Organização**
   - Coloque cada serviço em seu próprio arquivo (ex: `user_service.go`).
   - Mocks devem ficar em subpastas `mock_application/`.

9. **Proibição de Lógica de Infraestrutura**
   - Não inclua lógica de infraestrutura (ex: acesso a banco, autenticação, etc) diretamente nos serviços.
   - Delegue para repositórios, managers ou outros serviços.

10. **Exemplo Completo**

```go
type AuthService interface {
    Login(ctx context.Context, credentials domain.Credentials) (*domain.AuthSession, error)
}

type authService struct {
    userRepository   domain.UserRepository
    passwordManager  domain.PasswordManager
    authTokenManager domain.AuthTokenManager
    sessionDuration  time.Duration
}

func NewAuthService(sessionDuration time.Duration, userRepository domain.UserRepository, passwordManager domain.PasswordManager, authTokenManager domain.AuthTokenManager) AuthService {
    return &authService{
        sessionDuration:  sessionDuration,
        userRepository:   userRepository,
        passwordManager:  passwordManager,
        authTokenManager: authTokenManager,
    }
}

func (s *authService) Login(ctx context.Context, credentials domain.Credentials) (*domain.AuthSession, error) {
    if err := credentials.Validate(); err != nil {
        return nil, err
    }
    // ... restante da lógica ...
}
```

---

> Siga este padrão para garantir serviços coesos, testáveis e alinhados à arquitetura do projeto.
