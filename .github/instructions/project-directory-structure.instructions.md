---
description: Guia detalhado da estrutura de diretórios e arquivos do projeto mystery-gifter-api
applyTo: "**"
---
# Estrutura de Diretórios e Arquivos do Projeto `mystery-gifter-api`

Este documento descreve detalhadamente a estrutura de pastas e arquivos do projeto, destacando o propósito de cada diretório e os arquivos mais relevantes. Siga este guia para navegar, entender e manter a organização do repositório.

## Estrutura Geral

```
docker-compose.yml
go.mod
go.sum
tools.go
cmd/
internal/
pkg/
test/
```

### Arquivos na Raiz
- **docker-compose.yml**: Orquestração de containers para ambiente de desenvolvimento/testes.
- **go.mod / go.sum**: Gerenciamento de dependências Go Modules.
- **tools.go**: Dependências de ferramentas utilizadas no projeto.

---

## Diretórios Principais

### `cmd/`
- **api/main.go**: Ponto de entrada principal da aplicação (API REST).

### `internal/`
Organiza a lógica de domínio, aplicação e infraestrutura. Subdividido em:

#### `internal/application/`
- Serviços de aplicação (camada de orquestração de regras de negócio):
  - `auth_service.go`, `group_service.go`, `user_service.go`: Serviços principais.
  - `*_test.go`: Testes dos serviços.
  - `mock_application/`: Mocks para testes dos serviços de aplicação.

#### `internal/domain/`
- Entidades e regras de domínio:
  - `auth.go`, `group.go`, `user.go`, `identity.go`, `security.go`: Modelos e lógica de domínio.
  - `errors.go`: Definições de erros de domínio.
  - `*_test.go`: Testes das entidades e regras.
  - `build_domain/`: Builders para facilitar criação de entidades em testes.
  - `mock_domain/`: Mocks de interfaces de domínio para testes.

#### `internal/infra/`
- Implementações de infraestrutura:
  - `runner.go`: Inicialização e execução da aplicação.
  - `config/`: Configuração da aplicação (`config.go`, `config_test.go`).
  - `entrypoint/`: Camada de entrada REST:
    - `error_handler.go`, `middlewares.go`, `routes.go`: Infraestrutura de roteamento e tratamento de erros.
    - `rest/`: Controllers REST, DTOs e builders:
      - `auth_controller.go`, `group_controller.go`, `user_controller.go`: Controllers principais.
      - `*_dto.go`: Data Transfer Objects.
      - `*_test.go`: Testes dos controllers.
      - `build_rest/`: Builders para DTOs em testes.

#### `internal/outgoing/`
- Integrações externas:
  - `identity/`: Geradores de identidade (ex: UUID).
  - `postgres/`: Repositórios e modelos para persistência em PostgreSQL.
    - `migrations/`: Scripts de migração do banco de dados.
    - `build_postgres/`: Builders para testes.
    - `mock_postgres/`: Mocks para testes.
  - `security/`: Gerenciamento de autenticação e criptografia (ex: JWT, bcrypt).

### `pkg/`
- Pacotes utilitários reutilizáveis:
  - `validator/`: Validação de dados e testes.

### `test/`
- Utilitários e helpers para testes:
  - `helper/json.go`: Funções auxiliares para manipulação de JSON em testes.

---

## Recomendações de Navegação
- **Para lógica de negócio**: Explore `internal/domain/` e `internal/application/`.
- **Para endpoints REST**: Veja `internal/infra/entrypoint/rest/`.
- **Para integrações externas**: Consulte `internal/outgoing/`.
- **Para utilitários**: Use `pkg/` e `test/`.

## Boas Práticas
- Mantenha cada camada bem separada.
- Utilize os builders e mocks para facilitar testes.
- Siga a estrutura de pastas para novos módulos.

---

> Consulte este arquivo sempre que houver dúvidas sobre a localização de funcionalidades ou arquivos no projeto.
