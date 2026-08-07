# Aureon

[![CI](https://github.com/Santozz-x/Aureon/actions/workflows/ci.yml/badge.svg)](https://github.com/Santozz-x/Aureon/actions/workflows/ci.yml)

**Building the Infrastructure Behind Web3.**

Aureon é uma plataforma open source de infraestrutura Web3, multi-chain, API-first e AI-native. Veja o [CHARTER.md](CHARTER.md) para a visão completa, princípios e roadmap do projeto.

> Status: **v0.1.0** — núcleo (Wallet + Transaction API na ARC Network) funcionando, Fase 1 do [roadmap](docs/implementation-plan.md) em andamento. Ver [CHANGELOG.md](CHANGELOG.md).

## Estrutura do repositório

Monorepo Go organizado com [`go.work`](go.work), seguindo Clean Architecture (domain → usecase → adapter → infra) e injeção de dependência manual via construtores.

```
modules/
  contracts/         # chainport.Adapter — a porta blockchain-agnostic (sem dependências)
  platform/           # API Gateway — cmd/gateway + internal/{domain,usecase,adapter,infra}
  chains/
    arc/              # Adaptador da ARC Network (fase 1) implementando chainport.Adapter
  mcp/                # Servidor MCP — cmd/mcp-server (skeleton)
sdk/
  go/                 # SDK oficial em Go, cliente HTTP para a API da Aureon
```

Cada blockchain suportada é um módulo Go independente em `modules/chains/`, implementando a interface `chainport.Adapter` definida em `modules/contracts`. Isso preserva o princípio *Blockchain Agnostic*: o `platform` nunca depende diretamente de um chain específico, apenas da porta compartilhada.

## Documentação

| Documento | Conteúdo |
|---|---|
| [CHARTER.md](CHARTER.md) | Visão, missão, princípios e pilares da plataforma |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Decisões técnicas: camadas, monorepo, MCP como client, DI |
| [docs/implementation-plan.md](docs/implementation-plan.md) | Requisitos funcionais, épicos e tarefas por sprint |
| [docs/tradeoffs.md](docs/tradeoffs.md) | Decision journal — o que foi escolhido e o que foi sacrificado |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | Como contribuir |
| [CHANGELOG.md](CHANGELOG.md) | O que mudou em cada versão |

## Quick start

Requer Go 1.26+ e Docker (para o Postgres local).

```bash
go work sync
make db-up          # sobe um Postgres local via docker compose (dev only)

export AUREON_DATABASE_URL="postgres://aureon:aureon@localhost:5432/aureon?sslmode=disable"
export AUREON_KEYSTORE_ENCRYPTION_KEY="$(openssl rand -hex 32)"

make build          # compila todos os módulos
make vet            # go vet em todos os módulos
make run-gateway    # sobe o API Gateway em :8080 (aplica as migrações automaticamente)
make run-mcp        # sobe o skeleton do servidor MCP
```

O Gateway falha ao subir se `AUREON_DATABASE_URL` ou `AUREON_KEYSTORE_ENCRYPTION_KEY` (32 bytes em hex) não estiverem definidas — ambas guardam dados sensíveis (chaves privadas de carteira, criptografadas em repouso com AES-256-GCM; hashes de API key) e não têm valor padrão de propósito. Ver [docs/tradeoffs.md](docs/tradeoffs.md#tr-008-aureon-é-custodial-guarda-as-chaves-privadas-via-chainportkeystore-começando-com-storage-em-memória) para as ressalvas de segurança ainda pendentes antes de qualquer deploy real.

## Licença

Apache 2.0 — veja [LICENSE](LICENSE).
