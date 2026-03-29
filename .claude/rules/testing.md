---
paths:
  - "**/*_test.go"
---

# Convenções de Testes

- Nome da função: `Test_<Type>_<Method>`
- Cenários: `t.Run("should ... when ...", func(t *testing.T) { ... })`
- Descrições dos cenários em inglês
- Três fases obrigatórias com comentários: `// given`, `// when`, `// then`
- Um `gomock.Controller` por subtest (`mockCtrl := gomock.NewController(t)`)
- Builders de `build_domain/`, `build_postgres/`, `build_rest/` para criar objetos de teste
- `testify/assert` para todas as asserções; verificar sempre erro E resultado
- `context.Background()` para contextos em testes; `time.Duration` explícita para durações

**Obrigatório: sempre executar os testes após a implementação, identificar as falhas e corrigir.**

## Mocks e dependências

- Passar `nil` para dependências que não são usadas no cenário em vez de criar mocks desnecessários
- Usar `DoAndReturn` para validar o que está sendo passado para mocks que recebem dados importantes (ex: `Create`, `Update`)
- Usar `SetArg` do GoMock para popular o resultado de `GetContext`/`SelectContext` em vez de `DoAndReturn` com ponteiros

## Asserções obrigatórias por camada

**Serviços:**
- Em cenários de erro: sempre usar `assert.EqualError` para validar a mensagem exata do erro além de `assert.ErrorAs`
- Cobrir todos os cenários de falha de dependência (ex: quando `repository.Create` falha, quando `repository.Update` falha)

**Repositórios:**
- Em cenários de erro: usar `assert.ErrorContains` para validar o prefixo da mensagem de erro encapsulada (ex: `"error inserting group invite"`)
- Usar `SetArg(argIndex, value)` para popular o destino em chamadas de `GetContext` e `SelectContext`

**Controllers:**
- Em cenários de erro: além do status HTTP, decodificar o body como `entrypoint.WebError` e validar `Code` e `Message`

## Builders

- O nome da variável interna do builder deve ser o nome completo do tipo em camelCase (ex: `groupInviteDTO`, não `dto`; `groupSummary`, não `summary`)
