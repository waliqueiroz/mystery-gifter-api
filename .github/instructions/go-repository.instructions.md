---
description: Convenções para implementação de repositórios (Repository Pattern) no projeto mystery-gifter-api
applyTo: internal/infra/outgoing/postgres/**
---
# Convenções para Implementação de Repositórios (Repository Pattern)

Este documento define as diretrizes para implementação de repositórios responsáveis pela persistência e recuperação de entidades no projeto `mystery-gifter-api`, especialmente na camada `internal/infra/outgoing/postgres/`.

## 1. Estrutura e Nomeação

- Cada repositório deve ser implementado como um struct privado (ex: `userRepository`, `groupRepository`).
- O construtor público deve seguir o padrão `New<Entity>Repository(db DB) domain.<Entity>Repository`.
- O struct deve conter um campo `db` do tipo `DB` (interface de acesso ao banco).

## 2. Métodos CRUD

- Implemente métodos para operações básicas: `Create`, `GetByID`, `GetByEmail` (para usuários), `Update`, etc.
- Todos os métodos devem receber um `context.Context` como primeiro parâmetro.
- Métodos que retornam entidades devem retornar ponteiros para structs de domínio e erro: `(*domain.Entity, error)`.
- Métodos de escrita (`Create`, `Update`) devem retornar apenas `error`.

## 3. Queries e Execução

- Utilize o pacote `squirrel` para construção de queries SQL, sempre usando `PlaceholderFormat(squirrel.Dollar)`.
- Para inserções e atualizações, utilize `ExecContext`.
- Para buscas, utilize `GetContext` (para um único registro) ou `SelectContext` (para múltiplos registros).
- Sempre trate erros de queries, diferenciando entre `sql.ErrNoRows` (retorne erro de recurso não encontrado) e outros erros (retorne erro genérico com contexto).

## 4. Tratamento de Erros

- Para violações de unicidade (ex: email já cadastrado), utilize o tipo de erro do Postgres (`pq.Error`) e retorne um erro de domínio apropriado (ex: `domain.NewConflictError`).
- Para registros não encontrados, retorne `domain.NewResourceNotFoundError`.
- Sempre adicione contexto à mensagem de erro usando `fmt.Errorf("contexto: %w", err)`.

## 5. Transações

- Para operações que envolvem múltiplas queries dependentes (ex: criar grupo e usuários), utilize transações (`BeginTxx`, `Rollback`, `Commit`).
- Sempre utilize `defer tx.Rollback()` para garantir rollback em caso de erro.
- Só chame `Commit` explicitamente após todas as operações terem sucesso.

## 6. Mapeamento de Entidades

- Utilize funções de mapeamento dedicadas para converter structs do banco para structs de domínio (ex: `mapUserToDomain`, `mapGroupToDomain`).
- As funções de mapeamento devem ser privadas e seguir as convenções do arquivo [copilot-mapper-functions.instructions.md](.github/instructions/copilot-mapper-functions.instructions.md).

## 7. Organização e Localização

- Cada repositório deve estar em seu próprio arquivo, nomeado como `<entity>_repository.go`.
- Structs auxiliares para persistência (ex: `User`, `Group`) devem ser definidos no mesmo arquivo do repositório correspondente.

## 8. Logging

- Em caso de erro em operações críticas (ex: falha ao inserir), registre o erro usando `log.Println` antes de retornar.

## 9. Boas Práticas Gerais

- Nunca exponha detalhes de implementação do banco fora do repositório.
- Não inclua lógica de negócio nos repositórios; limite-se à persistência e recuperação de dados.
- Sempre use contextos para todas as operações de banco.
- Prefira retornar erros de domínio ao invés de erros genéricos.

## Exemplos

```go
// Estrutura do repositório
type userRepository struct {
	db DB
}

func NewUserRepository(db DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Exemplo de método Create
func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	query, args, err := squirrel.Insert("users").
		Columns("id", "name", ...).
		Values(user.ID, user.Name, ...).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building users insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var currentError *pq.Error
		if errors.As(err, &currentError) && currentError.Code.Name() == POSTGRES_UNIQUE_VIOLATION {
			return domain.NewConflictError("the email is already registered")
		}
		return fmt.Errorf("error inserting user: %w", err)
	}

	return nil
}

// Exemplo de método GetByID
func (r *userRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building users select query: %w", err)
	}

	var user User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return mapUserToDomain(user)
}
```

> Siga este padrão para garantir consistência, testabilidade e manutenibilidade dos repositórios em todo o projeto.
