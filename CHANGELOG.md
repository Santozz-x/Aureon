# Changelog

Todas as mudanças notáveis deste projeto são documentadas neste arquivo.

O formato segue [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), e o projeto adere a [Semantic Versioning](https://semver.org/lang/pt-BR/) — enquanto a versão principal for `0.x`, a API pode mudar sem aviso prévio.

## [Unreleased]

## [0.1.0] - 2026-08-06

Primeiro checkpoint público do projeto: Wallet API e Transaction API funcionando de ponta a ponta na ARC Network, com autenticação, persistência criptografada e CI. Cobre as Sprints 0–6 do [plano de implementação](docs/implementation-plan.md) — a Fase 1 (MVP núcleo) ainda não está completa (faltam OpenAPI, CLI e MCP, Sprints 7–10).

### Added

- Monorepo Go com `go.work` sobre 5 módulos independentes (`contracts`, `platform`, `chains/arc`, `mcp`, `sdk/go`), seguindo Clean Architecture e injeção de dependência manual.
- `chainport.Adapter`: porta blockchain-agnostic que todo módulo de chain implementa.
- API Gateway (`modules/platform`) com Wallet API e Transaction API reais na ARC Network:
  - `POST /v1/{network}/wallets` — cria carteira (custodial, chave gerada localmente).
  - `GET /v1/{network}/wallets/{address}/balance` — saldo real via `eth_getBalance`.
  - `POST /v1/{network}/transactions` — envia transação assinada (EIP-155) via `eth_sendRawTransaction`.
  - `POST /v1/{network}/transactions/estimate` — estimativa de gas via `eth_estimateGas`.
- Cliente RPC da ARC Network (`modules/chains/arc/rpc`), sobre `go-ethereum`, já que a Arc é EVM-compatível.
- Autenticação por API Key (`Authorization: Bearer aur_...`), com emissão/revogação (`POST /v1/apikeys`, `DELETE /v1/apikeys/{id}`) e middleware de autorização (`middleware.RequireAPIKey`).
- Persistência em PostgreSQL: chaves privadas de carteira criptografadas em repouso com AES-256-GCM (`infra/keystore.Postgres`), hashes de API key (`infra/apikeystore.Postgres`), migrações embutidas no binário via `go:embed` + `golang-migrate`.
- Logging estruturado (`middleware.Logging`) — método, path, status, duração e request ID em toda requisição.
- Suite de testes de integração ponta a ponta do Gateway (router + usecases + middlewares reais, só o adapter de chain é um fake).
- CI no GitHub Actions: build, vet, test (com Postgres real de serviço), `golangci-lint`, `gofmt -l`, em todos os módulos.
- SDK oficial em Go (`sdk/go`) — cliente HTTP mínimo para a API da Aureon.
- Skeleton do servidor MCP (`modules/mcp`) — ainda sem tools implementadas.
- Documentação: [CHARTER.md](CHARTER.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md), [docs/implementation-plan.md](docs/implementation-plan.md) (requisitos, épicos, sprints), [docs/tradeoffs.md](docs/tradeoffs.md) (decision journal), [docs/networks/arc.md](docs/networks/arc.md), [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md).

### Security

- Aureon é **custodial por design**: gera e guarda as chaves privadas das carteiras (nunca expostas ao chamador). Ver [TR-008](docs/tradeoffs.md#tr-008-aureon-é-custodial-guarda-as-chaves-privadas-via-chainportkeystore-começando-com-storage-em-memória).
- Chaves privadas são criptografadas em repouso (AES-256-GCM) antes de tocar o banco. A chave mestra de criptografia vem hoje de uma variável de ambiente (`AUREON_KEYSTORE_ENCRYPTION_KEY`), **não de um KMS** — ver [TR-010](docs/tradeoffs.md#tr-010-criptografia-em-repouso-com-chave-de-aplicação-aes-256-gcm-em-vez-de-kms). Não usar em produção com fundos reais antes disso mudar.
- `POST /v1/apikeys` está **intencionalmente público** nesta versão (sem fluxo de Identity/signup ainda) — qualquer um pode emitir uma API key para um `account_id` arbitrário. Ver [TR-009](docs/tradeoffs.md#tr-009-api-keys--hash-sha-256-não-bcryptargon2-storage-em-memória-emissão-sem-autenticação-prévia). **Bloqueador para qualquer deploy real.**

### Known limitations

- Apenas ARC Network (testnet), apenas Wallet + Transaction API — Explorer, Contract, Token e NFT APIs ainda não existem.
- Sem OpenAPI, CLI ou tools de MCP ainda (Sprints 7–10).
- Sem Dashboard, Playground, GraphQL, Webhooks/Event Streaming, suporte multi-chain além da ARC, ou infraestrutura de Cloud (Docker/K8s/Helm/Terraform) — todos backlog macro, ver [implementation-plan.md §7](docs/implementation-plan.md#7-backlog-macro--fases-2-4).
