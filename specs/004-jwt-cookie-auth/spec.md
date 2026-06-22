# Feature Specification: Autenticação Dual-Channel com Cookie HttpOnly

**Feature Branch**: `004-jwt-cookie-auth`
**Created**: 2026-06-21
**Status**: Draft

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Login com Sessão Segura para Web (Priority: P1)

Como usuário que acessa a aplicação pelo navegador, quero que minha sessão seja armazenada de forma que o JavaScript da página não consiga acessá-la, para que minha conta fique protegida contra ataques de roubo de credenciais via scripts maliciosos.

**Why this priority**: Resolve a principal vulnerabilidade de segurança atual (token no localStorage). É o núcleo da feature e habilita todas as demais histórias.

**Independent Test**: Pode ser testado realizando login e verificando que o navegador recebe um cookie de sessão marcado como inacessível via JavaScript, enquanto o aplicativo continua funcionando normalmente nas rotas protegidas.

**Acceptance Scenarios**:

1. **Given** usuário com credenciais válidas, **When** realiza login, **Then** recebe um cookie de sessão seguro no navegador e consegue acessar as rotas protegidas normalmente
2. **Given** usuário autenticado via cookie, **When** tenta acessar rota protegida, **Then** acesso é concedido sem necessidade de informar token manualmente
3. **Given** cookie de sessão presente, **When** código JavaScript da página tenta ler o cookie, **Then** o acesso é negado pelo navegador

---

### User Story 2 - Encerramento de Sessão Explícito (Priority: P2)

Como usuário autenticado via navegador, quero poder encerrar minha sessão de forma explícita, para que meu acesso seja revogado imediatamente e o cookie de sessão seja removido do navegador.

**Why this priority**: Complementa diretamente a autenticação via cookie — sem logout, o usuário não tem controle sobre o ciclo de vida da sessão.

**Independent Test**: Pode ser testado realizando logout e verificando que o cookie de sessão é removido e tentativas subsequentes de acessar rotas protegidas são rejeitadas.

**Acceptance Scenarios**:

1. **Given** usuário autenticado via cookie, **When** realiza logout, **Then** o cookie de sessão é removido do navegador
2. **Given** usuário que acabou de fazer logout, **When** tenta acessar rota protegida usando o cookie anterior, **Then** recebe resposta de não autorizado

---

### User Story 3 - Consulta de Dados do Usuário Logado (Priority: P3)

Como usuário autenticado, quero consultar meus próprios dados de perfil sem precisar saber ou informar meu identificador, para que eu possa obter minhas informações de forma simples e direta a qualquer momento.

**Why this priority**: Necessário para que o frontend possa exibir informações do usuário logado — especialmente após autenticação via cookie, onde o token não é acessível via JavaScript para extração de dados.

**Independent Test**: Pode ser testado realizando uma requisição ao endpoint de perfil com sessão válida e verificando o retorno das informações do usuário autenticado.

**Acceptance Scenarios**:

1. **Given** usuário autenticado via cookie, **When** consulta seus dados de perfil, **Then** recebe suas informações sem precisar informar nenhum identificador
2. **Given** usuário autenticado via cabeçalho de autorização, **When** consulta seus dados de perfil, **Then** recebe suas informações da mesma forma
3. **Given** usuário não autenticado, **When** tenta consultar dados de perfil, **Then** recebe resposta de não autorizado

---

### User Story 4 - Compatibilidade com Clientes Mobile (Priority: P4)

Como cliente mobile (Android ou iOS), quero continuar autenticando via cabeçalho HTTP de autorização, para que o comportamento atual dos apps nativos seja preservado sem necessidade de alterações.

**Why this priority**: Garantia de retrocompatibilidade. Aplicações mobile não utilizam cookies da mesma forma que navegadores — forçar cookies quebraria esses clientes.

**Independent Test**: Pode ser testado acessando qualquer rota protegida usando apenas o cabeçalho de autorização com token válido e verificando que o acesso é concedido normalmente, sem qualquer alteração no comportamento existente.

**Acceptance Scenarios**:

1. **Given** token de autenticação válido, **When** cliente envia apenas o cabeçalho de autorização (sem cookie), **Then** acesso à rota protegida é concedido
2. **Given** ambos cookie e cabeçalho de autorização válidos presentes, **When** cliente faz requisição, **Then** a autenticação é realizada via cookie (fonte prioritária) e o acesso é concedido normalmente

---

### Edge Cases

- O que acontece quando o cookie de sessão expirou e o usuário tenta acessar uma rota protegida?
- Como o sistema responde quando o cookie está presente mas contém um token corrompido ou adulterado?
- O que acontece quando o usuário tenta fazer logout sem estar autenticado?
- Quando cookie está presente mas inválido, o sistema rejeita a requisição (não faz fallback para o cabeçalho de autorização), pois o cookie tem precedência como fonte de autenticação
- O cookie deve funcionar em ambiente de desenvolvimento local (conexão sem HTTPS)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O sistema DEVE, ao realizar login com sucesso, definir automaticamente um cookie de sessão seguro no navegador do usuário, além de continuar retornando o token no corpo da resposta
- **FR-002**: O cookie de sessão DEVE ser configurado de forma a impedir acesso por scripts JavaScript da página (cookie httpOnly)
- **FR-003**: O cookie de sessão DEVE ser transmitido apenas em conexões seguras em ambiente de produção, e permitido em HTTP apenas em ambiente de desenvolvimento
- **FR-004**: O cookie de sessão DEVE ser configurado para prevenir envio automático em requisições originadas de domínios externos, protegendo contra ataques de falsificação de requisição entre sites (CSRF); como frontend e API operam em origens distintas (subdomínios ou portas diferentes), a política adotada deve permitir navegação legítima entre origens do mesmo domínio pai sem bloquear o fluxo de autenticação
- **FR-005**: O sistema DEVE proteger todas as rotas autenticadas aceitando tanto o cookie de sessão quanto o cabeçalho de autorização como formas válidas de autenticação; quando ambos estiverem presentes, o cookie tem precedência — o cabeçalho é avaliado apenas na ausência de cookie
- **FR-006**: O sistema DEVE disponibilizar um endpoint de logout que remova o cookie de sessão do navegador do usuário
- **FR-007**: O sistema DEVE disponibilizar o endpoint `GET /users/me` para consulta dos dados do usuário atualmente autenticado, retornando os mesmos campos públicos já expostos nos demais endpoints de usuário (excluindo credenciais internas), sem exigir que o identificador seja informado na requisição; por compartilhar o namespace `/users`, este endpoint deve ser declarado antes de `GET /users/:userID` para evitar conflito de roteamento
- **FR-008**: As responsabilidades de criação do cookie, validação da sessão e remoção do cookie DEVEM ser tratadas como operações distintas e isoladas, sem acoplamento ao fluxo de negócio de autenticação

### Key Entities

- **Sessão do Usuário**: Representa o estado de autenticação ativo de um usuário; contém identificador do usuário e período de validade
- **Cookie de Autenticação**: Mecanismo de persistência da sessão no navegador; associado ao domínio da aplicação e inacessível via JavaScript
- **Token de Autenticação**: Credencial de acesso portátil; pode ser transportado via cookie (para clientes web) ou cabeçalho HTTP (para clientes mobile)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Usuários web autenticados não têm o token de autenticação acessível por JavaScript da página, eliminando a exposição via localStorage
- **SC-002**: O fluxo completo de login, acesso a rotas protegidas e logout funciona sem que o usuário precise manipular manualmente tokens de autenticação no navegador
- **SC-003**: 100% dos clientes que já utilizam autenticação via cabeçalho de autorização continuam funcionando sem nenhuma modificação
- **SC-004**: O endpoint de consulta de dados do usuário logado retorna as informações corretamente para qualquer forma de autenticação (cookie ou cabeçalho)
- **SC-005**: 100% das rotas protegidas existentes aceitam ambos os mecanismos de autenticação sem regressão de comportamento

## Assumptions

- A aplicação opera tanto via HTTPS (produção) quanto HTTP (desenvolvimento local); o cookie deve funcionar em ambos os ambientes com configuração adequada por ambiente
- Frontend e API operam em origens distintas (subdomínios ou portas diferentes) em produção; a configuração do cookie deve ser compatível com esse cenário cross-origin
- Não há requisito de revogação de sessão no lado do servidor — a expiração natural do token é suficiente para o ciclo de vida da sessão
- **Não há mecanismo de refresh token nesta feature** — quando o token expira, o usuário deve realizar login novamente; renovação automática de sessão é fora de escopo
- O frontend web será atualizado para utilizar o novo mecanismo de cookie em substituição ao localStorage após a entrega desta feature
- Clientes mobile existentes não serão alterados e devem continuar funcionando exclusivamente via cabeçalho de autorização
- A duração do cookie de sessão deve ser equivalente à duração configurada atualmente para o token de autenticação
- Não há requisito de suporte a múltiplas sessões simultâneas por usuário em dispositivos distintos — uma sessão por dispositivo/browser é suficiente

## Clarifications

### Session 2026-06-22

- Q: O sistema deve renovar automaticamente a sessão antes da expiração, ou o usuário deve fazer login novamente quando o token expirar? → A: Sem refresh token — sessão expira e usuário faz login novamente
- Q: Quando uma requisição chega com cookie E cabeçalho de autorização, qual fonte é tentada primeiro? → A: Cookie tem precedência — cabeçalho só é avaliado na ausência de cookie
- Q: O que o endpoint de perfil deve retornar? → A: Dados completos do usuário (mesmo contrato dos endpoints existentes, excluindo credenciais internas); rota definida como `GET /users/me` para consistência com o namespace `/users` existente
- Q: Frontend web e API rodam na mesma origem (domínio + porta) em produção? → A: Não — origens diferentes (subdomínios ou portas distintas); política de cookie deve ser compatível com cross-origin entre domínio pai compartilhado
