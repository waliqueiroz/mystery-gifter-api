---
name: branch-conventions
description: |
  Define as convenções de nomeação de branches e direcionamento de PRs para features gerenciadas
  pelo speckit (padrão NNN-feature-name).

  INVOCAR APENAS dentro de um workflow speckit — ou seja, quando uma feature branch com o padrão
  NNN-feature-name (ex: 003-fiber-v3-upgrade, 005-mobile-ui-redesign) está ativa ou sendo discutida.
  Fora do speckit, usar git flow normal (feat/, fix/, refactor/) sem aplicar esta skill.

  Invocar proativamente quando TODAS as condições abaixo forem verdadeiras:
  - Uma feature branch speckit (NNN-feature-name) é o contexto atual
  - E uma das ações abaixo está prestes a acontecer:
    - Criar um branch com `git checkout -b` durante o `/speckit.implement`
    - Abrir um PR com `gh pr create` para um task/checkpoint da feature
    - O usuário pergunta para onde um PR deve apontar ou como nomear um branch dentro da feature
    - O usuário menciona "task branch", "checkpoint branch", "PR final da feature" ou "stacked PR"
    - O `/speckit.implement` está sendo executado (aplicar a cada branch e PR criado)

  NÃO invocar para: hotfixes, refactors, PRs avulsos ou qualquer trabalho de branch fora de
  uma feature speckit (esses seguem git flow com prefixos feat/, fix/, refactor/ normalmente).
---

# Convenções de Branches e PRs (Speckit)

> **Escopo**: estas convenções aplicam-se **exclusivamente** a features gerenciadas pelo speckit,
> identificadas pelo padrão `NNN-feature-name` (ex: `003-fiber-v3-upgrade`).
> Para hotfixes, refactors e qualquer outro trabalho fora do speckit, use git flow normal
> (`feat/`, `fix/`, `refactor/` etc.) sem aplicar as regras desta skill.

O projeto segue um modelo de branches em camadas dentro de cada feature speckit: task branches fluem para a feature branch, e só a feature branch flui para a `main`.

## Tipos de Branch

### Feature Branch
**Formato**: `NNN-feature-name`
**Criado por**: `/speckit.specify` automaticamente
**Exemplos**: `003-fiber-v3-upgrade`, `004-groups-profile-features`, `005-mobile-ui-redesign`

Nunca crie feature branches manualmente — o speckit cria e faz checkout automaticamente.

### Task / Checkpoint Branch
**Formato**: `task/NNN-{descrição-curta}`
**Criado por**: você, durante `/speckit.implement`
**Exemplos**: `task/003-cp1-go-update`, `task/003-cp2-fiber-v3`, `task/004-T001-T003-foundational`, `task/005-phase-3a-base-primitives`

Regras inegociáveis:
- Sempre criado A PARTIR da feature branch (ou do task branch anterior se stacking)
- PR sempre aponta para a feature branch (ou para o task anterior se stacking)
- **NUNCA aponta para `main` diretamente** — isso é o erro mais comum

## Criando um Task Branch

```bash
# 1. Garantir que está na feature branch (ou no task anterior)
git checkout 003-fiber-v3-upgrade

# 2. Criar o task branch
git checkout -b task/003-cp1-go-update

# ... implementar, commitar ...

# 3. Push e PR apontando para a FEATURE BRANCH
git push -u origin task/003-cp1-go-update
gh pr create \
  --base 003-fiber-v3-upgrade \
  --head task/003-cp1-go-update \
  --title "feat: descrição" \
  --body "..."
```

## Quando Usar Stacking entre Task Branches

Stack task B sobre task A (B aponta para A em vez da feature branch) quando:
- Task B **não compila ou não funciona** sem as mudanças de A — dependência de código real
- A dependência é **estritamente sequencial**

Exemplo real: atualizar o Go para 1.26.4 (cp1) antes de migrar o Fiber V3 (cp2) porque o Fiber V3 requer Go 1.25+.

```bash
# cp2 depende estritamente do cp1
git checkout task/003-cp1-go-update
git checkout -b task/003-cp2-fiber-v3

gh pr create \
  --base task/003-cp1-go-update \
  --head task/003-cp2-fiber-v3 \
  --title "feat: migrar Fiber V2 → V3" \
  --body "> **Stacked em**: #10 (task/003-cp1-go-update)"
```

Quando cp1 for mergeado, o GitHub automaticamente retargeta o PR de cp2 para a feature branch.

Se não há dependência real de código, sempre aponte direto para a feature branch — diffs menores, reviews mais limpos.

## PR Final: Feature Branch → main

Após TODOS os task PRs serem mergeados na feature branch:

```bash
# Verificar que a feature branch está limpa e os gates passam
git checkout 003-fiber-v3-upgrade
git pull

make build && make test   # ou equivalente do projeto

# Abrir o PR único que vai para main
gh pr create \
  --base main \
  --head 003-fiber-v3-upgrade \
  --title "feat(003): título descritivo da feature completa" \
  --body "$(cat <<'EOF'
## Resumo
- bullet 1
- bullet 2

## Tasks incluídas
- task/003-cp1-go-update (#10)
- task/003-cp2-fiber-v3 (#11)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

Este é o **único PR** que deve apontar para `main` em toda a lifetime de uma feature speckit.

## Diagrama

```
main
 └── 003-fiber-v3-upgrade              ← feature branch (speckit)  → PR final → main
      ├── task/003-cp1-go-update        → PR → 003-fiber-v3-upgrade
      └── task/003-cp2-fiber-v3         → PR → task/003-cp1-go-update (stacked)
                                               ↑ auto-retargeta para feature após cp1 merger
```

```
main
 └── 004-groups-profile-features       ← feature branch
      ├── task/004-T001-T003-foundational   → PR → 004-groups-profile-features
      ├── task/004-T004-T007-us1-cards      → PR → 004-groups-profile-features
      ├── task/004-T008-T013-us2-filtering  → PR → 004-groups-profile-features
      └── task/004-T014-T018-us3-profile    → PR → 004-groups-profile-features
```

## Checklist Antes de Qualquer PR

- [ ] Branch name começa com `task/NNN-`?
- [ ] `--base` é a feature branch (ou task anterior se stacking)?
- [ ] NÃO está apontando para `main`?
- [ ] Se é o PR final: todos os tasks mergeados, gates passando?
