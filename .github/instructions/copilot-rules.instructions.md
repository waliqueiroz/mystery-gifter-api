---
description: Como adicionar ou editar instruções do Copilot em seu projeto
---
# Guia de Gerenciamento de Regras do Copilot

## Formato da Estrutura de Regras

Toda regra do copilot deve seguir exatamente esta estrutura de metadados e conteúdo:

````markdown
---
description: Breve descrição do propósito da regra
applyTo: caminho/opcional/padrão/**/*
---
# Título da Regra

Conteúdo principal explicando a regra com formatação markdown.

1. Instruções passo a passo
2. Exemplos de código
3. Diretrizes

Exemplo:
```go
// Good
func goodExample() {
  // Correct implementation
}

// Bad example
func badExample() {
  // Incorrect implementation
}
```
````

## Organização de Arquivos

### Localização Obrigatória

Todos os arquivos de regras do copilot **devem** ser colocados em:

```
PROJECT_ROOT/.github/instructions/
```

### Estrutura de Diretórios

```
PROJECT_ROOT/
├── .github/
│   └── instructions/
│       ├── your-rule-name.instructions.md
│       ├── another-rule.instructions.md
│       └── one-more-rule.instructions.md
└── ...
```

### Convenções de Nomenclatura

- Use **kebab-case** para todos os nomes de arquivos
- Sempre use a extensão **.instructions.md**
- Faça nomes **descritivos** do propósito da regra
- Exemplos: `typescript-style.instructions.md`, `tailwind-styling.instructions.md`, `mdx-documentation.instructions.md`

## Diretrizes de Conteúdo

### Escrevendo Regras Eficazes

1. **Seja específico e acionável** - Forneça instruções claras
2. **Inclua exemplos de código** - Mostre boas e más práticas
3. **Faça referência a arquivos existentes** - Use links de arquivos md 
4. **Mantenha o foco** - Uma regra por preocupação/padrão
5. **Adicione contexto** - Explique por que a regra existe

### Formato dos Exemplos de Código

```typescript
// ✅ Bom: Claro e segue convenções
function processUser({ id, name }: { id: string; name: string }) {
  return { id, displayName: name };
}

// ❌ Ruim: Passagem de parâmetros pouco clara
function processUser(id: string, name: string) {
  return { id, displayName: name };
}
```

```go
// ✅ Bom: Nomenclatura clara, função pequena, usa struct para entrada/saída
type User struct {
	ID   string
	Name string
}

type ProcessedUser struct {
	ID          string
	DisplayName string
}

func ProcessUser(user User) ProcessedUser {
	return ProcessedUser{
		ID:          u.ID,
		DisplayName: u.Name,
	}
}

// ❌ Ruim: Nomenclatura vaga, retorna mapa anônimo, menos legível
func HandleData(id string, name string) map[string]string {
	return map[string]string{"id": id, "displayName": name}
}
```

### Referências a Arquivos

Ao referenciar arquivos do projeto nas regras, use links de arquivos md.

## Localizações Proibidas

**Nunca** coloque arquivos de regras em:
- Diretório raiz do projeto
- Qualquer subdiretório fora de `.github/instructions/`
- Diretórios de componentes
- Pastas de código-fonte
- Pastas de documentação

## Categorias de Regras

Organize as regras por propósito:
- **Estilo de Código**: `typescript-style.instructions.md`, `css-conventions.instructions.md`
- **Arquitetura**: `component-patterns.instructions.md`, `folder-structure.instructions.md`
- **Documentação**: `mdx-documentation.instructions.md`, `readme-format.instructions.md`
- **Ferramentas**: `testing-patterns.instructions.md`, `build-config.instructions.md`
- **Meta**: `copilot-rules.instructions.md`, `self-improve.instructions.md`

## Melhores Práticas

### Checklist de Criação de Regras
- [ ] Arquivo colocado no diretório `.github/instructions/`
- [ ] Nome do arquivo usa kebab-case com extensão `.instructions.md`
- [ ] Inclui seção de metadados adequada
- [ ] Contém título e seções claras
- [ ] Fornece exemplos bons e ruins
- [ ] Referencia arquivos relevantes do projeto
- [ ] Segue formatação consistente

### Manutenção
- **Revise regularmente** - Mantenha as regras atualizadas com mudanças na base de código
- **Atualize exemplos** - Garanta que os exemplos de código reflitam os padrões atuais
- **Faça referências cruzadas** - Conecte regras relacionadas
- **Documente alterações** - Atualize as regras quando os padrões evoluírem