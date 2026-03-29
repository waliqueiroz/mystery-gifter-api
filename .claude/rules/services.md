---
paths:
  - "internal/application/**/*.go"
---

# Convenções da Camada de Serviço

- Pacote deve ser `application`
- `//go:generate go run go.uber.org/mock/mockgen -destination mock_application/x_service.go . XService` deve ser a primeira linha após `package application`
- Interface pública `XService`, struct privada `xService`, construtor `NewXService(...) XService`
- Construtor deve retornar a interface (nunca `*xService`)
- Todos os métodos: `context.Context` como primeiro parâmetro, `error` como último retorno
- Chamar `entity.Validate()` antes de qualquer lógica de negócio
- Mocks gerados em `mock_application/`
- Regras de negócio que envolvem verificação de permissão ou estado de uma entidade devem ser encapsuladas em métodos da própria entidade (domínio rico) — o serviço apenas chama o método e propaga o erro
