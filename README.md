# Aureon

**Building the Infrastructure Behind Web3.**

Aureon é uma plataforma open source de infraestrutura Web3, multi-chain, API-first e AI-native. Veja o [CHARTER.md](CHARTER.md) para a visão completa, princípios e roadmap do projeto.

> Status: **Draft** — em bootstrap inicial.

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

## Quick start

Requer Go 1.26+.

```bash
go work sync

make build          # compila todos os módulos
make vet            # go vet em todos os módulos
make run-gateway    # sobe o API Gateway em :8080
make run-mcp        # sobe o skeleton do servidor MCP
```

## Licença

Apache 2.0 — veja [LICENSE](LICENSE).
