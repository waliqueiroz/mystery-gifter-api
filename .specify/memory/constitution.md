<!--
## Relatório de Impacto de Sincronização

**Mudança de Versão**: Constituição do Mystery Gifter API 1.0.0 → 1.1.0

**Princípios Modificados**: N/A

**Seções Adicionadas**:
- Princípio VIII: Idioma dos Artefatos (INEGOCIÁVEL) — todos os artefatos do speckit em pt-BR

**Seções Removidas**: N/A

**Templates que Requerem Atualização**:
- ✅ `.specify/templates/spec-template.md` — Nota de idioma adicionada no cabeçalho
- ✅ `.claude/commands/speckit.specify.md` — Instrução explícita de idioma adicionada

**TODOs de Acompanhamento**: Nenhum — todos os campos resolvidos.
-->

# Constituição da Mystery Gifter API

## Princípios Fundamentais

### I. Clean Architecture (INEGOCIÁVEL)

O código DEVE manter estrita separação em três camadas: **domain**, **application** e
**infrastructure**. A direção de dependência flui apenas para dentro — infrastructure depende
de application, application depende de domain, domain não depende de nada externo.

- `internal/domain/` DEVE conter apenas entidades, interfaces, erros de domínio e validação.
  Nenhuma importação de infraestrutura é permitida.
- `internal/application/` DEVE orquestrar a lógica de domínio apenas via interfaces injetadas.
  Sem chamadas diretas a banco de dados, HTTP ou frameworks.
- `internal/infra/` DEVE implementar as interfaces de domain/application.
  Lógica de negócio NÃO DEVE aparecer em controllers ou repositórios.
- Todo serviço público DEVE seguir: interface pública (`XService`), struct privada (`xService`),
  construtor retornando a interface (`NewXService(...) XService`).
- Todo repositório DEVE seguir: struct privada (`xRepository`), construtor retornando
  `domain.XRepository` (`NewXRepository(db DB) domain.XRepository`).

**Justificativa**: A violação de limites de camada cria código não testável e acoplamento forte com
escolhas de infraestrutura. Esta regra foi estabelecida desde o início do projeto e nunca é negociável.

### II. Disciplina de Testes

Toda lógica de negócio DEVE ter testes unitários escritos no mesmo ciclo de implementação.
Os testes DEVEM ser executados e passar antes que uma feature seja considerada completa.

- Nomenclatura de funções de teste: `Test_<Tipo>_<Método>`
- Subtestes DEVEM usar `t.Run("should ... when ...", ...)` com descrições em inglês
- Todo teste DEVE ter três fases explícitas anotadas com comentários: `// given`, `// when`, `// then`
- Cada subteste DEVE criar seu próprio `gomock.Controller`: `mockCtrl := gomock.NewController(t)`
- Dados de teste DEVEM ser construídos via builders em `build_domain/`, `build_postgres/`,
  ou `build_rest/` — nunca structs literais inline para objetos complexos
- Todas as asserções DEVEM usar `testify/assert`; tanto erro QUANTO resultado DEVEM ser sempre verificados
- `context.Background()` DEVE ser usado para contextos em testes
- **Obrigatório**: executar testes após cada implementação; identificar e corrigir falhas antes de prosseguir

**Justificativa**: Testes escritos depois são incompletos e perdem casos extremos. O padrão
given/when/then garante rastreabilidade entre cenários de teste e critérios de aceite.

### III. Validação Orientada ao Domínio

Entidades e value objects DEVEM se auto-validar. A lógica de validação pertence à camada de
domain, não a controllers ou services.

- Toda entidade DEVE implementar um método `Validate() error` usando tags `go-playground/validator`
- Funções de fábrica (`NewX(...)`) DEVEM chamar `Validate()` antes de retornar
- Métodos mutadores (`AddUser`, `Archive`, etc.) DEVEM chamar `Validate()` após a mutação
- Erros de domínio DEVEM usar os tipos customizados definidos em `internal/domain/errors.go`:
  `ValidationError`, `ConflictError`, `ResourceNotFoundError`, `UnauthorizedError`, `ForbiddenError`
- Cada tipo de erro DEVE carregar o código de status HTTP correto como parte de sua definição
- Controllers NÃO DEVEM definir mensagens de erro de negócio — eles propagam erros de domínio diretamente

**Justificativa**: Centralizar a validação no domínio garante consistência em todos os pontos de
entrada (REST, CLI, consumidores de fila) e torna os invariantes de domínio explícitos e testáveis.

### IV. Contrato de API Consistente

Todos os controllers REST DEVEM seguir um fluxo de requisição-resposta estrito e uniforme para
garantir comportamento previsível ao cliente.

- Assinatura de handler: `func (c *XController) Method(ctx *fiber.Ctx) error`
- Fluxo de requisição: `BodyParser` → `dto.Validate()` → `mapXToDomain` → chamada ao service →
  `mapXFromDomain` → resposta
- Falha em `BodyParser` DEVE retornar `fiber.NewError(fiber.StatusUnprocessableEntity)`
- Todos os outros erros DEVEM ser retornados diretamente (o error handler mapeia erros de domínio para status HTTP)
- Criação de recurso DEVE retornar `ctx.Status(fiber.StatusCreated).JSON(...)`
- Todas as outras respostas bem-sucedidas DEVEM retornar `ctx.JSON(...)`
- Parâmetros de rota DEVEM usar `ctx.Params("paramName")`, nunca `ctx.Query` para IDs de recursos
- ID do usuário autenticado DEVE ser extraído via `c.AuthTokenManager.GetAuthUserID(ctx.Locals("user"))`
- Anotações Swagger DEVEM ser mantidas para todos os endpoints; `make generate-docs` DEVE passar
- Parâmetros de query no Swagger DEVEM ser definidos individualmente (nunca `schema: "$ref"`)

**Justificativa**: Fluxo uniforme de controller reduz carga cognitiva e previne inconsistências
no tratamento de erros que vazam detalhes internos para os clientes.

### V. Abstração de Infraestrutura

Todas as dependências externas DEVEM ser acessadas através de interfaces, nunca por tipos concretos.

- Repositórios DEVEM injetar a interface `DB`, nunca `*sqlx.DB` diretamente
- Todas as queries SQL DEVEM usar `squirrel` com `PlaceholderFormat(squirrel.Dollar)`
- Operações de escrita DEVEM usar `ExecContext`, leituras de linha única `GetContext`,
  leituras de múltiplas linhas `SelectContext`
- `sql.ErrNoRows` DEVE ser mapeado para `domain.NewResourceNotFoundError`
- Violação de unicidade PostgreSQL (`pq.Error`) DEVE ser mapeada para `domain.NewConflictError`
- Todos os outros erros DEVEM ser encapsulados: `fmt.Errorf("mensagem de contexto: %w", err)`
- Escritas em múltiplos passos DEVEM usar transações: `BeginTxx` → `defer tx.Rollback()` →
  operações → `tx.Commit()`
- Mocks DEVEM ser gerados com `go.uber.org/mock/mockgen` via diretivas `//go:generate`;
  mudanças de interface requerem reexecução de `go generate ./...`

**Justificativa**: Design orientado a interfaces permite testes unitários isolados e desacopla a
lógica de negócio de escolhas específicas de infraestrutura (engine de banco, backend de storage).

### VI. Simplicidade & YAGNI

Toda abstração DEVE justificar sua existência resolvendo um problema presente, não um hipotético
futuro. Complexidade requer justificativa explícita.

- Sem helpers, utilitários ou abstrações para operações únicas
- Sem features especulativas, configuração extra ou shims de compatibilidade retroativa
- Sem docstrings, comentários ou anotações de tipo adicionados a código que não foi alterado
- Tratamento de erro NÃO DEVE ser adicionado para cenários que não podem ocorrer dados os invariantes atuais
- Funções mapeadoras DEVEM ser privadas e co-localizadas com o tipo que servem
- Slices DEVEM ser pré-alocados com `make([]T, 0, len(src))` antes da iteração
- Três linhas de código semelhantes são preferíveis a uma abstração prematura

**Justificativa**: Abstrações prematuras aumentam o custo de manutenção e obscurecem a intenção. A
quantidade certa de complexidade é exatamente o que a tarefa requer.

### VII. Performance & Observabilidade

A API DEVE permanecer responsiva sob a carga esperada. Erros DEVEM ser observáveis sem
expor detalhes internos aos clientes.

- Todos os métodos de repositório DEVEM passar `context.Context` para habilitar propagação de timeout
- `log.Println` DEVE ser chamado para erros significativos de infraestrutura antes de retorná-los
- Paginação DEVE ser aplicada a todos os endpoints de listagem/busca via filtros `Limit` + `Offset`
- Valores padrão de paginação DEVEM ser definidos como constantes na camada de domínio
  (ex.: `DefaultGroupLimit = 15`)
- Direção e campo de ordenação DEVEM ser validados via tags `oneof` do `go-playground/validator`
- Tempos de resposta HTTP DEVEM ser adequados para uso interativo (meta p95 < 200ms para
  CRUD padrão; nenhum mecanismo de enforcement de SLO obrigatório neste estágio)

**Justificativa**: Propagação de contexto e logging estruturado de erros são a linha de base mínima
de observabilidade. Paginação previne queries sem limite de degradar o banco de dados.

### VIII. Idioma dos Artefatos (INEGOCIÁVEL)

Existe uma separação rígida de idioma entre linguagem humana e linguagem de máquina:

**Português brasileiro (pt-BR) OBRIGATÓRIO em**:
- Todos os artefatos do speckit: `spec.md`, `plan.md`, `tasks.md`, checklists, relatórios de análise
- Títulos de seções, descrições, critérios de aceite, requisitos funcionais e critérios de sucesso
- Comentários em código-fonte (explicações do porquê, não do quê)
- Documentação Swagger (`summary`, `description` das anotações)
- Mensagens de commit e descrições de pull request
- Respostas e explicações da IA assistente

**Inglês OBRIGATÓRIO em**:
- Todo código-fonte: nomes de funções, variáveis, tipos, structs, constantes, pacotes
- Mensagens de erro retornadas pelo sistema (strings em `errors.New`, `fmt.Errorf`, `fiber.NewError`)
- Descrições de subtestes (`t.Run("should ... when ...", ...)`) conforme o Princípio II
- Nomes de branch, identificadores técnicos (FR-001, SC-001, P1, etc.), caminhos de arquivo
- Placeholders de template (ex.: `[FEATURE NAME]`, `[DATE]`)
- Tags de validação, nomes de campos JSON/query, headers HTTP

**Justificativa**: O time de desenvolvimento opera em português brasileiro, mas código em inglês
é padrão universal na indústria e facilita onboarding, uso de bibliotecas e colaboração open source.
Misturar idiomas em código cria inconsistência e dificulta manutenção. A separação clara elimina
ambiguidades sobre onde cada idioma se aplica.

## Padrões Tecnológicos

**Linguagem**: Go (versão estável mais recente)
**Framework Web**: Fiber v2 — usar apenas padrões idiomáticos do Fiber
**Banco de Dados**: PostgreSQL via `sqlx` + query builder `squirrel`
**Migrações**: `golang-migrate` — arquivos de migração DEVEM ser commitados junto com mudanças de schema
**Autenticação**: `golang-jwt` — tokens JWT; chave secreta e duração da sessão via variáveis de ambiente
**Validação**: `go-playground/validator` — tags de struct são a fonte canônica de validação
**Testes**: `testify/assert` para asserções, `go.uber.org/mock/mockgen` para mocks
**Documentação**: `go-swagger` — `make generate-docs` DEVE ter sucesso sempre
**Identidade**: UUID v4 via `internal/infra/outgoing/identity` (injetado como `IdentityGenerator`)
**Configuração**: `caarlos0/env` — toda configuração via variáveis de ambiente; `.env.example` DEVE
ser mantido atualizado

Variáveis de ambiente obrigatórias: `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`,
`DB_PASSWORD`, `AUTH_SECRET_KEY`, `AUTH_SESSION_DURATION`.

## Workflow de Desenvolvimento

- Para arquivos com mais de 100 linhas ou mudanças complexas: esboçar o plano e confirmar a abordagem antes de editar
- Perspectivas alternativas DEVEM ser levantadas quando houver oportunidade de melhorar a qualidade
- `context.Context` DEVE ser o primeiro parâmetro de todo método de service e repositório
- `error` DEVE ser o último valor de retorno de todo método que pode falhar
- Após toda mudança de interface, executar `go generate ./...` para regenerar mocks
- Após toda mudança de DTO ou rota, executar `make generate-docs` para atualizar a spec Swagger
- Todos os testes unitários DEVEM passar antes de marcar qualquer tarefa como concluída (`make test`)
- Mensagens de commit DEVEM ser descritivas e referenciar a camada alterada
  (ex.: `feat(domain): add group archiving logic`)

**Gates de Qualidade** (todos DEVEM passar antes de uma feature ser considerada concluída):

1. `make build` — compila sem erros
2. `make test` — todos os testes unitários passam
3. `make generate-docs` — spec Swagger gerada sem erros
4. Nenhum token de placeholder inexplicado permanece em nenhum documento de spec ou plano

## Governança

Esta constituição substitui todas as outras práticas e diretrizes de código do projeto. Qualquer regra
no `CLAUDE.md` ou em `.claude/rules/` que conflite com esta constituição DEVE ser reconciliada
em favor deste documento, e a regra conflitante atualizada de acordo.

**Procedimento de emenda**:
1. Propor mudança com justificativa na descrição do pull request
2. Atualizar `CONSTITUTION_VERSION` conforme versionamento semântico:
   - MAJOR: remoção de princípio, redefinição ou mudança de governança incompatível com versão anterior
   - MINOR: novo princípio ou seção adicionado, ou orientação materialmente expandida
   - PATCH: esclarecimento, correção de texto ou refinamento não-semântico
3. Atualizar `LAST_AMENDED_DATE` para a data da emenda (formato ISO YYYY-MM-DD)
4. Executar o checklist de propagação de consistência contra todos os arquivos em `.specify/templates/`
5. Documentar o impacto no comentário do Relatório de Impacto de Sincronização no topo deste arquivo

**Revisão de conformidade**: Todo plano de feature DEVE incluir um gate "Verificação da Constituição"
confirmando que o design proposto não viola nenhum princípio. Violações DEVEM ser
justificadas em uma tabela de Rastreamento de Complexidade no documento de plano.

Orientação de desenvolvimento em tempo real: ver `CLAUDE.md` e `.claude/rules/` para convenções
específicas da linguagem que complementam estes princípios.

**Versão**: 1.1.0 | **Ratificada**: 2026-03-28 | **Última Emenda**: 2026-06-21
