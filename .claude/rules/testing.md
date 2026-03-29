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
