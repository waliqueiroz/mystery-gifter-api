---
description: 
applyTo: "**"
---
# DIRETRIZES OPERACIONAIS DE EDIÇÃO
                
## DIRETIVA PRINCIPAL
	Evite trabalhar em mais de um arquivo por vez.
	Edições simultâneas múltiplas em um arquivo causarão corrupção.
	Converse e ensine sobre o que você está fazendo enquanto codifica.

## PROTOCOLO PARA ARQUIVOS GRANDES E MUDANÇAS COMPLEXAS

### FASE DE PLANEJAMENTO OBRIGATÓRIA
	Ao trabalhar com arquivos grandes (>300 linhas) ou mudanças complexas:
		1. SEMPRE comece criando um plano detalhado ANTES de fazer quaisquer edições
            2. Seu plano DEVE incluir:
                   - Todas as funções/seções que precisam de modificação
                   - A ordem em que as mudanças devem ser aplicadas
                   - Dependências entre as mudanças
                   - Número estimado de edições separadas necessárias
                
            3. Formate seu plano como:
## PLANO DE EDIÇÃO PROPOSTO
	Trabalhando com: [nome do arquivo]
	Total de edições planejadas: [número]

### FAZENDO EDIÇÕES
	- Foque em uma mudança conceitual por vez
	- Mostre trechos claros do "antes" e "depois" ao propor mudanças
	- Inclua explicações concisas do que mudou e por quê
	- Sempre verifique se a edição mantém o estilo de código do projeto

### Sequência de edição:
	1. [Primeira mudança específica] - Propósito: [por quê]
	2. [Segunda mudança específica] - Propósito: [por quê]
	3. Você aprova este plano? Prosseguirei com a Edição [número] após sua confirmação.
	4. AGUARDE confirmação explícita do usuário antes de fazer QUAISQUER edições quando o usuário ok editar [número]
            
### FASE DE EXECUÇÃO
	- Após cada edição individual, indique claramente o progresso:
		"✅ Completada edição [#] de [total]. Pronto para próxima edição?"
	- Se você descobrir mudanças adicionais necessárias durante a edição:
	- PARE e atualize o plano
	- Obtenha aprovação antes de continuar
                
### ORIENTAÇÃO PARA REFATORAÇÃO
	Ao refatorar arquivos grandes:
	- Divida o trabalho em partes logicamente independentes e funcionais
	- Garanta que cada estado intermediário mantém a funcionalidade
	- Considere duplicação temporária como uma etapa intermediária válida
	- Sempre indique o padrão de refatoração sendo aplicado
                
### EVITANDO RATE LIMIT
	- Para arquivos muito grandes, sugira dividir as mudanças em várias sessões
	- Priorize mudanças que são unidades logicamente completas
	- Sempre forneça pontos claros de parada
            
## Requisitos Gerais
	Use tecnologias modernas para todas as sugestões de código. Priorize código limpo e manutenível, evitando comentários em excesso.
                             
## Considerações de Segurança
	- Sanitize todas as entradas do usuário minuciosamente.
	- Parametrize consultas de banco de dados.
	- Aplique Políticas de Segurança de Conteúdo (CSP) fortes.
	- Use proteção CSRF onde aplicável.
	- Garanta cookies seguros (`HttpOnly`, `Secure`, `SameSite=Strict`).
	- Limite privilégios e aplique controle de acesso baseado em função.
	- Implemente logging e monitoramento interno detalhado.