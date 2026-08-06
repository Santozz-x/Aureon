# Plano de Implementação

Converte o [CHARTER.md](../CHARTER.md) em requisitos funcionais, épicos e tarefas de sprint. Ver decisões e trade-offs em [tradeoffs.md](tradeoffs.md) e a arquitetura técnica em [ARCHITECTURE.md](ARCHITECTURE.md).

## Premissas

- **Equipe:** 1 desenvolvedor (solo).
- **Sprint:** 1 semana.
- **Capacidade:** ~8 story points/sprint (escala Fibonacci: 1, 2, 3, 5, 8, 13). Ver [TR-006](tradeoffs.md#tr-006-capacidade-assumida-de-8-story-points-por-sprint-solo-1-semana).
- **Horizonte detalhado:** Fase 1 (MVP núcleo — ARC Network) está quebrada em sprints com tarefas acionáveis. Fases 2+ cobrem os demais pilares do charter como Epics com estimativa macro, a detalhar progressivamente. Ver [TR-005](tradeoffs.md#tr-005-horizonte-de-planejamento--todos-os-pilares-documentados-como-epics-mas-só-a-fase-1-mvp-arc-detalhada-em-sprints).
- **Estado atual:** o bootstrap do monorepo (Sprint 0) já está commitado — ver histórico do git.

---

## 1. Requisitos Funcionais

Derivados dos Platform Pillars do charter.

### Infrastructure (INF)

| ID | Requisito |
|---|---|
| FR-001 | API Gateway expõe endpoints REST unificados por rede (`/v1/{network}/...`) |
| FR-002 | Wallet API: criar carteira, consultar saldo |
| FR-003 | Transaction API: enviar transação, estimar gas, consultar status |
| FR-004 | Explorer API: buscar transações, blocos e endereços |
| FR-005 | Contract API: deploy e leitura/execução de contratos |
| FR-006 | Token API: transferências e consultas de tokens fungíveis |
| FR-007 | NFT API: mint, transferência e consulta de NFTs |
| FR-008 | Identity: modelo de conta/organização de desenvolvedor |
| FR-009 | Authentication: autenticação de requisições (API Key, e futuramente OAuth) |
| FR-010 | Authorization: controle de acesso por escopo/role |
| FR-011 | API Keys: emissão, rotação e revogação |
| FR-012 | Webhooks: notificação de eventos on-chain via callback HTTP |
| FR-013 | Event Streaming: stream de eventos em tempo real (WebSocket/SSE) |

### Developer Experience (DX)

| ID | Requisito |
|---|---|
| FR-014 | Dashboard web para gestão de API keys e uso |
| FR-015 | CLI (`aureon`) para operações via terminal |
| FR-016 | Playground interativo para testar chamadas de API |
| FR-017 | Especificação OpenAPI publicada e versionada |
| FR-018 | API GraphQL como alternativa ao REST |
| FR-019 | Documentação (guias, referência de API, exemplos) |
| FR-020 | Templates de projeto (starter kits) |
| FR-021 | Quickstarts por linguagem/framework |

### Artificial Intelligence (AI)

| ID | Requisito |
|---|---|
| FR-022 | Servidor MCP expondo: `createWallet`, `transfer`, `deployContract`, `readContract`, `mintNFT`, `searchTransactions`, `estimateGas`, `bridgeAssets`, `swapTokens`, `listenEvents` |

### Cloud Platform (CLOUD)

| ID | Requisito |
|---|---|
| FR-023 | Imagens Docker para todos os serviços |
| FR-024 | Manifests Kubernetes |
| FR-025 | Helm charts |
| FR-026 | Módulos Terraform para provisionamento |
| FR-027 | Suporte multi-região |
| FR-028 | Auto scaling |

### Multi-Chain (CHAIN)

| ID | Requisito |
|---|---|
| FR-029 | Adapter ARC Network (Fase 1, já em andamento) |
| FR-030 | Adapters adicionais via módulos independentes (Ethereum, Base, Arbitrum, Optimism, Polygon, BNB Chain, Solana, Cosmos, Avalanche, Bitcoin) |

### SDKs (SDK)

| ID | Requisito |
|---|---|
| FR-031 | SDK oficial Go |
| FR-032 | SDKs oficiais em outras linguagens (TypeScript, Python, ...) |

---

## 2. Epics

| Epic | Nome | FRs | Esforço macro (SP) | Fase |
|---|---|---|---|---|
| E1 | Core Platform & Auth | FR-001, FR-008, FR-009, FR-010, FR-011 | 26 | 1 |
| E2 | Wallet API (ARC) | FR-002, FR-029 | 16 | 1 |
| E3 | Transaction API (ARC) | FR-003 | 16 | 1 |
| E4 | Persistência & Observabilidade | — (transversal) | 16 | 1 |
| E5 | DX básica (OpenAPI + CLI MVP) | FR-015, FR-017 | 16 | 1 |
| E6 | MCP mínimo (3 tools) | FR-022 (parcial) | 16 | 1 |
| E7 | Explorer API | FR-004 | 21 | 2 |
| E8 | Contract API | FR-005 | 34 | 2 |
| E9 | Token API | FR-006 | 21 | 2 |
| E10 | NFT API | FR-007 | 21 | 2 |
| E11 | Webhooks & Event Streaming | FR-012, FR-013 | 34 | 3 |
| E12 | Dashboard | FR-014 | 55 | 3 |
| E13 | Playground, GraphQL, Templates, Quickstarts | FR-016, FR-018, FR-020, FR-021 | 55 | 3 |
| E14 | Documentação completa | FR-019 | 21 | 3 (contínuo) |
| E15 | MCP completo (tools restantes) | FR-022 (restante) | 34 | 3 |
| E16 | Cloud Platform | FR-023–FR-028 | 55 | 4 |
| E17 | Multi-chain expansion | FR-030 | 34+/chain | 4 |
| E18 | SDKs multi-linguagem | FR-032 | 21+/linguagem | 4 |

Fase 1 = MVP núcleo detalhado em sprints abaixo. Fases 2–4 = backlog macro, ordenado por dependência natural (Explorer/Contract/Token/NFT dependem do mesmo padrão de Wallet/Transaction já estabelecido; Dashboard/Playground dependem da API estar estável; Cloud e multi-chain expansion vêm depois de haver algo em produção para operar).

---

## 3. Estrutura de arquivos (Fase 1)

Arquivos/diretórios já existentes (bootstrap) marcados ✅; a criar marcados ⬜.

| Caminho | Épico | Status |
|---|---|---|
| `modules/contracts/chain.go` | E2/E3 | ✅ |
| `modules/platform/internal/usecase/wallet.go` | E2 | ✅ |
| `modules/platform/internal/adapter/rest/{handler,router}.go` | E1/E2 | ✅ |
| `modules/platform/internal/infra/config/config.go` | E1 | ✅ |
| `modules/platform/cmd/gateway/main.go` | E1 | ✅ |
| `modules/chains/arc/adapter.go` | E2/E3 | ✅ (stubs, sem RPC real) |
| `modules/chains/arc/rpc/client.go` | E2/E3 | ✅ Sprint 1 |
| `modules/chains/arc/rpc/client_test.go` | E2/E3 | ✅ Sprint 1 |
| `docs/networks/arc.md` | E2/E3 | ✅ Sprint 1 |
| `modules/contracts/chain.go` (`KeyStore`) | E2 | ✅ Sprint 2 |
| `modules/platform/internal/infra/keystore/memory.go` | E2 | ✅ Sprint 2 (interino — ver TR-008) |
| `modules/platform/internal/usecase/transaction.go` | E3 | ✅ Sprint 3 |
| `modules/platform/internal/usecase/registry.go` (resolver de adapter compartilhado) | E2/E3 | ✅ Sprint 3 |
| `modules/platform/internal/adapter/rest/transaction.go` | E3 | ✅ Sprint 3 |
| `modules/chains/arc/rpc/client.go` (`PendingNonceAt`, `SuggestGasPrice`) | E3 | ✅ Sprint 3 |
| `modules/platform/internal/domain/identity.go` | E1 | ✅ Sprint 4 |
| `modules/platform/internal/adapter/rest/middleware/auth.go` | E1 | ✅ Sprint 4 |
| `modules/platform/internal/usecase/apikey.go` | E1 | ✅ Sprint 4 |
| `modules/platform/internal/infra/apikeystore/memory.go` | E1 | ✅ Sprint 4 (interino — ver TR-009) |
| `modules/platform/internal/adapter/rest/apikey.go` | E1 | ✅ Sprint 4 |
| `modules/platform/internal/infra/keystore/postgres.go` | E4 | ✅ Sprint 5 (Postgres real, não `adapter/repository/` como planejado originalmente — ver nota abaixo) |
| `modules/platform/internal/infra/apikeystore/postgres.go` | E4 | ✅ Sprint 5 |
| `modules/platform/internal/infra/db/db.go` (conexão + runner de migrações) | E4 | ✅ Sprint 5 |
| `modules/platform/migrations/` (`.sql` + `go:embed`) | E4 | ✅ Sprint 5 |
| `docker-compose.yml` (Postgres local, dev only) | E4 | ✅ Sprint 5 |
| `.github/workflows/ci.yml` | E4 | ⬜ Sprint 6 |
| `api/openapi.yaml` | E5 | ⬜ Sprint 7 |
| `modules/cli/` (novo módulo, `cmd/aureon`) | E5 | ⬜ Sprint 8 |
| `modules/mcp/internal/tools/` | E6 | ⬜ Sprint 9–10 |

---

## 4. Sprints — Fase 1 (MVP núcleo)

### Sprint 0 — Bootstrap (concluído)

Monorepo `go.work`, módulo `contracts` (porta `chainport.Adapter`), `platform` com Gateway HTTP mínimo (`/health`, criação/consulta de wallet via use case), `chains/arc` com adapter stub (métodos retornam "not implemented"), `mcp` skeleton, `sdk/go` com client HTTP mínimo. Build, vet e smoke test validados. Charter, README, LICENSE e este conjunto de docs commitados.

### Sprint 1 — Cliente RPC da ARC Network (8 SP) ✅ concluído

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-101 | Pesquisar e documentar a API RPC/HTTP da ARC Network (endpoints, autenticação, formatos) | E2 | FR-029 | 2 | — | ✅ Nota técnica em [`docs/networks/arc.md`](networks/arc.md): Arc é EVM-compatível (JSON-RPC `eth_*`), chain ID `5042002` (testnet), gas nativo em USDC, RPC configurável via `AUREON_ARC_RPC_URL` |
| T-102 | Implementar cliente HTTP/RPC genérico em `modules/chains/arc/rpc` | E2/E3 | FR-029 | 3 | T-101 | ✅ `Client` (via `go-ethereum`/`ethclient`) expõe `ChainID`, `BalanceAt`, `EstimateGas`, `SendTransaction`; erros do transporte/upstream sempre encapsulados com `%w` |
| T-103 | Testes do cliente RPC com `httptest` (mock do servidor ARC) | E2/E3 | FR-029 | 3 | T-102 | ✅ `go test ./...` cobre `ChainID`, `BalanceAt` (sucesso e erro upstream), `EstimateGas` (sucesso e erro upstream) |

Decisão de dependência (go-ethereum, geração local de chave) registrada em [TR-007](tradeoffs.md#tr-007-cliente-rpc-da-arc-via-go-ethereum-e-geração-local-de-chave-em-vez-de-circle-wallets).

### Sprint 2 — Wallet API real (8 SP) ✅ concluído

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-104 | Implementar `CreateWallet` real no Adapter ARC (geração de keypair/endereço) | E2 | FR-002 | 5 | T-102 | ✅ `POST /v1/arc/wallets` retorna endereço válido (testado contra a testnet real: `0xFbB49CAfEDdFDbFCd5F6FDF6E5D289eAc3b7055B`); chave persistida via `chainport.KeyStore` |
| T-105 | Implementar `GetBalance` real via cliente RPC | E2 | FR-002 | 3 | T-102 | ✅ `GET /v1/arc/wallets/{address}/balance` retorna saldo real (testado contra a testnet: endereço de burn `0x000...dead` → `6550183353001544053029` USDC, unidade menor) |

Decisão de custódia (KeyStore custodial, em memória por ora) registrada em [TR-008](tradeoffs.md#tr-008-aureon-é-custodial-guarda-as-chaves-privadas-via-chainportkeystore-começando-com-storage-em-memória). `chainport.Adapter` ganhou `KeyStore` como porta adicional — atualizar a nota de arquitetura se necessário.

### Sprint 3 — Transaction API real (8 SP) ✅ concluído

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-106 | Implementar `SendTransaction` no Adapter ARC | E3 | FR-003 | 5 | T-104 | ✅ Monta, assina (EIP-155/legacy tx via go-ethereum) e envia via `eth_sendRawTransaction`; testado contra a testnet real — pipeline completo (nonce, gas price, chain id, assinatura, broadcast) confirmado, rejeitado corretamente por saldo zero (carteira nova sem fundos de testnet) |
| T-107 | Implementar `EstimateGas` no Adapter ARC | E3 | FR-003 | 3 | T-102 | ✅ Testado contra a testnet real (mesma ressalva de saldo zero acima) e com testes unitários (mock RPC) |
| — | Expor endpoints `POST /v1/{network}/transactions` e `POST /v1/{network}/transactions/estimate` no `usecase`/`rest` | E3 | FR-001, FR-003 | (incluso acima) | T-106, T-107 | ✅ `usecase.TransactionService` + `rest.TransactionHandler`, roteados no Gateway |

**Nota para a próxima sessão de testes manuais:** para validar um `SendTransaction` com sucesso (não só a rejeição por saldo), a carteira `from` precisa de USDC de testnet — usar o faucet oficial da Arc (ver [docs/networks/arc.md](networks/arc.md)) antes de repetir o smoke test.

### Sprint 4 — Identity, Auth e API Keys (8 SP) ✅ concluído

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-108 | Modelar domínio de Identity (conta de desenvolvedor) e API Key (`internal/domain/identity.go`) | E1 | FR-008, FR-011 | 2 | — | ✅ `domain.Account`, `domain.APIKey` (guarda só `SecretHash`, nunca o segredo em texto plano) |
| T-109 | Middleware de autenticação por API Key no Gateway | E1 | FR-009 | 3 | T-108 | ✅ `middleware.RequireAPIKey` — validado manualmente: sem header → `401`, header inválido → `401`, chave válida → `200` e chega no handler |
| T-110 | Endpoint de emissão/revogação de API Keys (storage em memória por ora) | E1 | FR-010, FR-011 | 3 | T-108 | ✅ `POST /v1/apikeys` cria (retorna o segredo uma única vez), `DELETE /v1/apikeys/{id}` revoga; validado manualmente que a chave revogada falha no middleware (`401`) |

Decisões registradas em [TR-009](tradeoffs.md#tr-009-api-keys--hash-sha-256-não-bcryptargon2-storage-em-memória-emissão-sem-autenticação-prévia): hash SHA-256 (correto para tokens de alta entropia, não senha), storage em memória (mesmo padrão do TR-008), e `POST /v1/apikeys` propositalmente público por enquanto — **isso precisa mudar antes de qualquer deploy real**.

### Sprint 5 — Persistência (8 SP) ✅ concluído

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-111 | Definir interfaces de repositório (`Store`) para Identity/API Keys, desacopladas de storage concreto | E4 | — | 2 | T-108 | ✅ `usecase.APIKeyStore` (já existia desde a Sprint 4) e `chainport.KeyStore` — ambas sem dependência de driver de DB |
| T-112 | Implementar repositório Postgres + migrações | E4 | — | 5 | T-111 | ✅ `infra/keystore.Postgres` (AES-256-GCM em repouso) e `infra/apikeystore.Postgres`; migrações em `modules/platform/migrations/` (embutidas via `go:embed`, aplicadas automaticamente no boot via `golang-migrate`); testes de integração contra Postgres real (`TEST_DATABASE_URL`, gated/skip se ausente) |
| T-113 | Configuração de conexão via env (`AUREON_DATABASE_URL`) | E4 | — | 1 | T-112 | ✅ `AUREON_DATABASE_URL` e `AUREON_KEYSTORE_ENCRYPTION_KEY` obrigatórias — Gateway falha rápido com mensagem clara se ausentes/inválidas (testado) |

Implementações em memória (`keystore.Memory`, `apikeystore.Memory`) mantidas apenas para testes unitários — o Gateway não as usa mais. Decisão da chave de criptografia via env var (não KMS) registrada em [TR-010](tradeoffs.md#tr-010-criptografia-em-repouso-com-chave-de-aplicação-aes-256-gcm-em-vez-de-kms), com ressalva explícita: **não usar em deploy com fundos reais** até virar KMS.

**Nota de estrutura:** o plano original previa um pacote genérico `internal/adapter/repository/postgres/`. Na implementação, ficou mais idiomático manter `Postgres` no mesmo pacote que seu par `Memory` (`infra/keystore/{memory,postgres}.go`, `infra/apikeystore/{memory,postgres}.go`) — mesma interface, implementações alternativas lado a lado, sem introduzir uma camada genérica extra.

**Validado end-to-end:** wallet criada e API key emitida num processo do Gateway sobrevivem a matar o processo e subir um binário novo — saldo e `estimateGas` funcionaram no segundo processo usando a mesma chave privada recuperada (descriptografada) do Postgres.

### Sprint 6 — Observabilidade, testes de integração e CI (8 SP)

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-114 | Logging estruturado consistente (`slog`) em todas as camadas, com request ID | E4 | — | 2 | Sprint 4/5 | Toda requisição HTTP loga método, path, status, duração e request ID |
| T-115 | Testes de integração do Gateway ponta a ponta (`httptest`, sem mocks internos) | E4 | — | 3 | Sprint 1–5 | Suite cobre wallet, transaction e auth com cenários de sucesso e erro |
| T-116 | CI no GitHub Actions: build, vet, test e lint para todos os módulos do `go.work` | E4 | — | 3 | T-115 | Pipeline verde em push/PR; falha se `golangci-lint` ou testes quebrarem |

### Sprint 7 — OpenAPI (8 SP)

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-117 | Especificação OpenAPI 3 cobrindo todos os endpoints existentes | E5 | FR-017 | 5 | Sprint 1–6 | `api/openapi.yaml` válido (lint via `spectral` ou similar), reflete request/response reais |
| T-118 | Servir a spec e uma UI de documentação (Swagger UI/Redoc) a partir do Gateway | E5 | FR-017, FR-019 | 3 | T-117 | `GET /docs` renderiza a documentação interativa |

### Sprint 8 — CLI MVP (8 SP)

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-119 | Scaffold do módulo `modules/cli` (`cmd/aureon`, Cobra), consumindo `sdk/go` | E5 | FR-015 | 3 | Sprint 1–3 | `aureon --help` funciona; binário compila via `make build` |
| T-120 | Comando `aureon wallet create --network arc` | E5 | FR-015 | 3 | T-119 | Cria wallet real via API e imprime o endereço |
| T-121 | Comando `aureon wallet balance --network arc --address ...` | E5 | FR-015 | 2 | T-119 | Imprime saldo real consultado via API |

### Sprint 9 — MCP mínimo, parte 1 (8 SP)

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-122 | Avaliar e integrar um SDK MCP para Go (transporte stdio) | E6 | FR-022 | 2 | — | `modules/mcp` inicializa um servidor MCP funcional (handshake completo com um client MCP de teste) |
| T-123 | Implementar tool `createWallet` no MCP, usando `sdk/go` como client do Gateway (ver [TR-004](tradeoffs.md#tr-004-mcp-server-como-cliente-http-da-api-via-sdkgo-não-como-segunda-implementação-de-negócio)) | E6 | FR-022 | 3 | T-122, Sprint 2 | Um agente MCP de teste consegue chamar `createWallet` e receber um endereço real |
| T-124 | Implementar tool `estimateGas` | E6 | FR-022 | 3 | T-122, Sprint 3 | Agente MCP recebe estimativa de gas real |

### Sprint 10 — MCP mínimo, parte 2 (8 SP)

| ID | Tarefa | Epic | FR | SP | Depende de | Critério de aceite |
|---|---|---|---|---|---|---|
| T-125 | Implementar tool `transfer` | E6 | FR-022 | 5 | T-122, Sprint 3, Sprint 4 (auth) | Agente MCP envia uma transação real via tool e recebe o `TxHash` |
| T-126 | Testes de integração das tools MCP (client MCP de teste automatizado) | E6 | FR-022 | 3 | T-123–T-125 | Suite automatizada cobre as 3 tools, incluindo casos de erro (rede não suportada, saldo insuficiente) |

**Critério de saída da Fase 1:** Sprints 0–10 concluídos ⇒ Wallet API e Transaction API funcionando de ponta a ponta na ARC Network, com auth, persistência, observabilidade, CI, OpenAPI, CLI e 3 tools MCP. Nesse ponto, detalhar a Fase 2 (Explorer/Contract/Token/NFT APIs) em sprints, atualizando este documento.

---

## 5. Caminho crítico (Fase 1)

```
T-101 → T-102 → T-103
           ↓
         T-104 → T-105
           ↓
         T-106 → T-107
           ↓
  T-108 → T-109 → T-110
           ↓
  T-111 → T-112 → T-113
           ↓
  T-114 → T-115 → T-116
           ↓
  T-117 → T-118
           ↓
  T-119 → T-120 → T-121
           ↓
  T-122 → T-123 → T-124 → T-125 → T-126
```

O cliente RPC da ARC (T-102) é o bloqueador central: Wallet, Transaction e, por consequência, CLI e MCP dependem dele. Priorizar Sprint 1 sem desvios.

---

## 6. Riscos

| Risco | Impacto | Probabilidade | Mitigação |
|---|---|---|---|
| API/RPC da ARC Network mal documentada ou instável | Alto — bloqueia E2/E3 inteiros | Média | T-101 como spike isolado antes de comprometer sprints seguintes; validar contra testnet cedo |
| Geração/custódia de chaves privadas implementada incorretamente (Wallet API) | Alto — risco de segurança/perda de fundos | Média | Revisão de segurança dedicada antes de qualquer deploy além de testnet; nunca logar chaves privadas |
| Capacidade real do solo dev divergir de 8 SP/sprint | Médio — cronograma impreciso | Alta | Reestimar após Sprint 1–2 (ver [TR-006](tradeoffs.md#tr-006-capacidade-assumida-de-8-story-points-por-sprint-solo-1-semana)) |
| SDK MCP para Go imaturo ou inexistente | Médio — atrasa E6 | Média | T-122 inclui avaliação explícita antes de comprometer T-123–126; alternativa é implementar o transporte stdio manualmente |
| Escopo "plataforma completa" (Fases 2–4) mudar antes de a Fase 1 terminar | Baixo — replanejamento | Alta (esperado) | Backlog de Epics em vez de tarefas fixas para Fases 2–4; revisar a cada fim de fase |

---

## 7. Backlog macro — Fases 2–4

Não quebrado em sprints ainda (ver [TR-005](tradeoffs.md#tr-005-horizonte-de-planejamento--todos-os-pilares-documentados-como-epics-mas-só-a-fase-1-mvp-arc-detalhada-em-sprints)). Ordem sugerida por dependência natural:

**Fase 2 — Completar as APIs de infraestrutura (ARC Network)**
E7 Explorer API → E8 Contract API → E9 Token API → E10 NFT API. Cada uma segue o mesmo padrão arquitetural de E2/E3 (porta em `contracts`, implementação no `chains/arc`, use case + handler REST no `platform`).

**Fase 3 — Experiência de desenvolvedor e IA completas**
E11 Webhooks & Event Streaming, E12 Dashboard, E13 Playground/GraphQL/Templates/Quickstarts, E14 Documentação completa (contínuo desde a Fase 1), E15 MCP completo (`deployContract`, `readContract`, `mintNFT`, `searchTransactions`, `bridgeAssets`, `swapTokens`, `listenEvents` — cada um depende da API correspondente já existir).

**Fase 4 — Escala: nuvem e multi-chain**
E16 Cloud Platform (Docker → Kubernetes → Helm → Terraform → multi-região → autoscaling, nessa ordem), E17 Multi-chain expansion (um novo módulo `modules/chains/{network}` por rede, reaproveitando 100% de `contracts` e `platform`), E18 SDKs multi-linguagem (TypeScript primeiro, por ser o mais demandado por devs web3).
