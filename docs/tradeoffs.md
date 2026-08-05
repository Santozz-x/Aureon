# Decision Journal (Tradeoffs)

Registra decisões de arquitetura e escopo que tinham alternativas reais, e o que foi sacrificado em cada escolha. Ver [ARCHITECTURE.md](ARCHITECTURE.md) para o estado atual e [implementation-plan.md](implementation-plan.md) para como isso vira tarefas.

<!-- ENTRY TEMPLATE
## TR-NNN: Título curto da decisão
- **Fase:** plan | build | review
- **Decisão:** o que foi escolhido
- **Alternativas consideradas:** o que mais foi avaliado
- **Sacrifício:** o que essa escolha custa ou deixa pior
- **Reavaliar quando:** condição objetiva que justificaria revisitar
-->

## TR-001: Clean Architecture + monorepo modular (go.work)
- **Fase:** plan
- **Decisão:** Organizar `platform` em camadas domain/usecase/adapter/infra, com cada blockchain como módulo Go independente unido via `go.work`.
- **Alternativas consideradas:** estrutura flat (um único `main.go` crescendo organicamente); um único módulo Go para todo o repositório.
- **Sacrifício:** mais boilerplate e mais arquivos `go.mod`/wiring desde o dia 1, para um projeto que ainda não tem nenhum usuário real.
- **Reavaliar quando:** nunca, a menos que o princípio "Blockchain Agnostic" seja abandonado — a estrutura existe especificamente para sustentar esse princípio do charter.

## TR-002: Injeção de dependência manual (sem framework de DI)
- **Fase:** plan
- **Decisão:** Wiring via construtores explícitos em `cmd/*/main.go`, sem `samber/do`, `wire` ou `dig/fx`.
- **Alternativas consideradas:** `google/wire` (geração de código), `samber/do` (container em runtime).
- **Sacrifício:** conforme o número de serviços/dependências crescer, o wiring manual em `main.go` pode ficar repetitivo.
- **Reavaliar quando:** um `main.go` único ultrapassar ~10 dependências wired manualmente, ou quando houver múltiplos serviços com grafos de dependência muito similares (sinal de duplicação).

## TR-003: `contracts` como módulo Go próprio (não `internal/`)
- **Fase:** plan
- **Decisão:** A porta `chainport.Adapter` vive em `modules/contracts`, um módulo público e independente.
- **Alternativas consideradas:** definir a interface dentro de `internal/domain` do `platform` e reexportá-la.
- **Sacrifício:** nenhum sacrifício real identificado — é a única opção que permite módulos de chain externos implementarem a interface sem acoplar a `platform`.
- **Reavaliar quando:** não aplicável.

## TR-004: MCP server como cliente HTTP da API (via `sdk/go`), não como segunda implementação de negócio
- **Fase:** plan
- **Decisão:** `modules/mcp` vai chamar a API do Gateway através do `sdk/go`, em vez de importar `usecase` diretamente.
- **Alternativas consideradas:** dar ao `mcp` acesso direto ao `usecase.WalletService` (exigiria tornar esse pacote público, fora de `internal/`).
- **Sacrifício:** uma chamada de rede a mais (MCP → Gateway) em vez de uma chamada em processo; latência levemente maior.
- **Reavaliar quando:** se o overhead de rede se provar um problema real de performance para agentes de IA (medir antes de otimizar).

## TR-005: Horizonte de planejamento — todos os pilares documentados como Epics, mas só a Fase 1 (MVP ARC) detalhada em sprints
- **Fase:** plan
- **Decisão:** O usuário pediu planejamento para "a plataforma completa" (todos os pilares do charter). Documentamos **todos** os Requisitos Funcionais e Epics correspondentes a todos os pilares, mas só quebramos em tarefas de sprint a Fase 1 (núcleo: Wallet + Transaction API na ARC Network, Auth, DX básica, MCP mínimo). As fases seguintes ficam como Epics com estimativa macro, a detalhar progressivamente.
- **Alternativas consideradas:** gerar tarefas de sprint detalhadas para todos os ~13 Epics imediatamente.
- **Sacrifício:** o roadmap de longo prazo (Explorer/Contract/Token/NFT APIs, Dashboard, Cloud, multi-chain, SDKs multi-linguagem) não tem tarefas acionáveis ainda — só estimativa de esforço em Epic.
- **Reavaliar quando:** a Fase 1 estiver perto da conclusão (ver critério de saída em [implementation-plan.md](implementation-plan.md)) — nesse ponto, detalhar a Fase 2 em sprints.

## TR-006: Capacidade assumida de 8 story points por sprint (solo, 1 semana)
- **Fase:** plan
- **Decisão:** Sprints de 1 semana dimensionados para ~8 SP, assumindo 1 desenvolvedor trabalhando neste projeto part-time/mixed com outras responsabilidades.
- **Alternativas consideradas:** nenhuma capacidade explícita informada pelo usuário; poderia ter assumido tempo integral (~13-20 SP/semana).
- **Sacrifício:** se a capacidade real for maior, os sprints vão parecer subdimensionados e podem ser combinados; se for menor, vão estourar.
- **Reavaliar quando:** após os 2 primeiros sprints reais executados — ajustar a capacidade com base na velocidade observada.

## TR-007: Cliente RPC da ARC via go-ethereum e geração local de chave, em vez de Circle Wallets
- **Fase:** build
- **Decisão:** `modules/chains/arc/rpc` usa `github.com/ethereum/go-ethereum` (`ethclient`, `core/types`, pacote `ethereum`) para falar com o JSON-RPC padrão da Arc, e `CreateWallet` vai gerar um par de chaves secp256k1 localmente em vez de usar o produto Circle Wallets (custódia gerenciada pela Circle, recomendado na doc oficial da Arc).
- **Alternativas consideradas:** cliente JSON-RPC feito à mão (`net/http` + `encoding/json`) sem depender de `go-ethereum`; usar Circle Wallets para custódia de chaves.
- **Sacrifício:** `go-ethereum` é uma dependência pesada (dezenas de pacotes transitivos, incluindo criptografia e otel) só para falar HTTP JSON-RPC e assinar transações. Em troca, evita reimplementar assinatura EIP-155/RLP à mão — risco alto de acertar errado em código que mexe com chaves privadas. Não usar Circle Wallets significa que o Aureon é responsável por proteger as chaves geradas (custódia própria), não a Circle.
- **Reavaliar quando:** se o tamanho do binário/build da `chains/arc` se tornar um problema real; ou se, ao adicionar a segunda chain EVM, ficar claro que vale a pena extrair um pacote `chains/evmshared` para não duplicar o wrapper do go-ethereum por módulo.
