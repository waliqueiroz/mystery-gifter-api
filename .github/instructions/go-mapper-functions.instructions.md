---
description: Convenções para funções de mapeamento (mapper functions) entre structs (DTO, domínio, persistência) no projeto mystery-gifter-api
applyTo: "**"
---
# Convenções para Funções de Mapeamento (Mapper Functions)

Padronize todas as funções de mapeamento entre structs (DTOs, domínio, modelos de persistência) seguindo estas diretrizes:

## 1. Nomeação

- Use o padrão `map<Entity>To<Destino>` para conversão direta (ex: banco → domínio).
- Use o padrão `map<Entity>From<Origem>` para conversão reversa (ex: domínio → DTO).
- Para coleções, use o plural: `mapUsersToDomain`, `mapUsersFromDomain`.

## 2. Estrutura da Função

- Receba o struct de origem como parâmetro (valor ou slice).
- Crie o struct de destino explicitamente, campo a campo.
- Sempre que o struct de destino possuir método `Validate()`, invoque-o antes de retornar.
- Se a validação falhar, retorne imediatamente o erro.
- Retorne sempre um ponteiro para o struct de destino e um erro (`(*Tipo, error)`), ou um slice de valores e erro (`[]Tipo, error`).

## 3. Mapeamento de Slices

- Para slices, itere sobre cada elemento, mapeando individualmente.
- Se algum elemento falhar na validação, retorne o erro imediatamente.
- Acumule os resultados em um slice e retorne.

## 4. Encadeamento

- Se o struct de destino possuir campos compostos (ex: outros structs), utilize funções de mapeamento auxiliares para esses campos.

## 5. Localização

- Mantenha as funções de mapeamento próximas dos tipos que convertem (ex: no mesmo arquivo do DTO, model de banco, etc).

## 6. Simplicidade

- Não inclua lógica de negócio nas funções de mapeamento.
- Limite-se à transformação de dados e validação.

## 7. Exemplo

```go
func mapUserToDomain(user User) (*domain.User, error) {
	domainUser := domain.User{
		ID:        user.ID,
		Name:      user.Name,
		// ... outros campos ...
	}
	if err := domainUser.Validate(); err != nil {
		return nil, err
	}
	return &domainUser, nil
}

func mapUsersToDomain(users []User) ([]domain.User, error) {
	domainUsers := make([]domain.User, 0, len(users))
	for _, model := range users {
		user, err := mapUserToDomain(model)
		if err != nil {
			return nil, err
		}
		domainUsers = append(domainUsers, *user)
	}
	return domainUsers, nil
}
```

> Siga este padrão para garantir consistência, legibilidade e facilidade de manutenção em todo o projeto.
