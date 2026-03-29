---
paths:
  - "internal/infra/outgoing/postgres/**/*.go"
---

# Convenções de Repositórios

- Struct privada `xRepository`, construtor `NewXRepository(db DB) domain.XRepository`
- Injetar interface `DB` (nunca `*sqlx.DB` diretamente — facilita testes e abstração)
- `squirrel` com `PlaceholderFormat(squirrel.Dollar)` para todas as queries SQL
- `ExecContext` para writes, `GetContext` para single row, `SelectContext` para múltiplas linhas
- `sql.ErrNoRows` → `domain.NewResourceNotFoundError`
- `pq.Error` com unique violation → `domain.NewConflictError`
- Encapsular outros erros com `fmt.Errorf("mensagem: %w", err)`
- Transações multi-step: `BeginTxx` → `defer tx.Rollback()` → operações → `tx.Commit()`
- Logging com `log.Println` para erros importantes antes de retorná-los
- Mocks gerados em `mock_postgres/`
