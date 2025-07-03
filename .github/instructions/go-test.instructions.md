---
description: Convenções e boas práticas para testes unitários em Go no projeto mystery-gifter-api
applyTo: internal/**/*_test.go
---
# Convenções para Testes Unitários em Go

Este guia define o padrão para a escrita de testes unitários no projeto, promovendo clareza, isolamento, legibilidade e cobertura adequada.

## Estrutura dos Testes

1. **Organização por função/método**  
   - Cada função/método público deve ter um bloco `t.Run` para cada cenário relevante.
   - O nome do teste segue o padrão: `Test_<Tipo>_<Método>` ou `Test_<funcionalidade>`.
   - Use subtestes (`t.Run`) para cobrir diferentes casos de uso e erros.

2. **Fases do Teste**  
   - Separe claramente as fases:  
     - `// given` (preparação/mocks/entrada)  
     - `// when` (execução da ação)  
     - `// then` (asserções/resultados esperados)  
   - Use comentários para marcar cada fase.

3. **Mocks e Dependências**  
   - Utilize o `gomock` para mocks de interfaces.
   - Instancie o `gomock.Controller` em cada subteste para garantir isolamento.
   - Use builders (ex: `build_domain.NewUserBuilder()`) para criar entidades de teste.

4. **Asserções**  
   - Use o pacote `testify/assert` para todas as asserções.
   - Prefira `assert.NoError`, `assert.Error`, `assert.Equal`, `assert.Nil`, `assert.ErrorIs` para clareza.
   - Sempre valide tanto o erro quanto o resultado retornado.

5. **Cobertura de Casos**  
   - Inclua casos de sucesso, falha de validação, falha de dependências, erros inesperados e edge cases.
   - Exemplos de cenários:  
     - Sucesso na operação  
     - Erro de autenticação/autorização  
     - Falha de integração (ex: token, hash, banco)  
     - Dados inválidos ou ausentes

6. **Nomenclatura dos Testes**  
   - O nome do subteste (`t.Run`) deve descrever claramente o cenário, ex:  
     - "should return status 201 and the user when the user is created successfully"
     - "should return an error when token creation fails"

7. **Isolamento**  
   - Cada subteste deve ser independente, sem dependência de estado global.
   - Sempre crie novos mocks e dados para cada subteste.

8. **Estilo e Organização**  
   - Imports organizados em blocos padrão Go.
   - Evite comentários excessivos, foque em nomes claros e autoexplicativos.
   - Não misture múltiplas responsabilidades em um único teste.

## Exemplos

```go
func Test_authService_Login(t *testing.T) {
	t.Run("should return a session successfully", func(t *testing.T) {
		// given
		// ...preparação de mocks e dados...

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.NoError(t, err)
		assert.Equal(t, authSession, *result)
	})

	t.Run("should return an error when token creation fails", func(t *testing.T) {
		// given
		// ...preparação de mocks e dados...

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
```

## Diretrizes Adicionais

- Use builders e mocks do projeto para facilitar a manutenção dos testes.
- Prefira nomes descritivos para variáveis de entrada e saída.
- Sempre cubra casos de erro e edge cases relevantes.
- Mantenha os testes rápidos e determinísticos.

> Atualize esta regra sempre que novos padrões de teste surgirem ou forem adotados no projeto.
