---
name: generate-docs
description: Gera a documentação Swagger/OpenAPI do projeto Mystery Gifter API. Use quando precisar atualizar ou visualizar a documentação da API após mudanças em DTOs ou endpoints.
user-invocable: true
argument-hint: "[serve]"
---

# Gerar Documentação Swagger

Execute `make generate-docs` para gerar a especificação em `docs/specs/swagger.yaml`.

Se $ARGUMENTS contiver "serve", execute também `make serve-docs` para servir a Swagger UI em http://localhost:8081.
