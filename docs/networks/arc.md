# ARC Network (Arc, da Circle)

Nota técnica de referência para `modules/chains/arc`. Entrega da tarefa T-101 ([implementation-plan.md](../implementation-plan.md)).

## O que é

Arc é uma blockchain Layer-1 desenvolvida pela Circle, "purpose-built for programmable money". Características relevantes para o adapter:

- **EVM-compatível**: expõe a API JSON-RPC padrão do Ethereum (`eth_*`) — mesmo protocolo usado por qualquer node Ethereum. Ferramentas padrão (Hardhat, Foundry, Viem, Ethers, go-ethereum) funcionam sem adaptação.
- **Gas nativo em USDC**: ao contrário de outras chains EVM, o token usado para pagar gas é o próprio USDC, não um token nativo separado.
- **Finalidade determinística sub-segundo**, via consenso Malachite (Tendermint BFT) + camada de execução separada.
- **Apenas testnet até o momento** desta nota (agosto/2026). Mainnet prevista para o verão de 2026 (hemisfério norte), segundo o whitepaper da Circle.

## Rede e RPC

| Item | Valor |
|---|---|
| Chain ID (testnet) | `5042002` |
| Moeda nativa (gas) | USDC |
| RPC HTTP (Circle, oficial) | `https://rpc.testnet.arc.io` |
| RPC WebSocket (Circle, oficial) | `wss://rpc.testnet.arc.io` |
| RPC HTTP (Blockdaemon) | `https://rpc.blockdaemon.testnet.arc.io` |
| RPC HTTP/WS (dRPC) | `https://rpc.drpc.testnet.arc.io` |
| RPC HTTP/WS (QuickNode) | `https://rpc.quicknode.testnet.arc.io` |

A documentação oficial não especifica se os endpoints exigem API key — tratar como **configurável por ambiente** (`AUREON_ARC_RPC_URL`), já que qualquer um dos providers acima (ou um node próprio) pode ser usado sem mudança de código. Mainnet ainda não tem endpoints publicados; quando existirem, a mesma env var muda de valor — nenhuma mudança de código é necessária.

## Métodos JSON-RPC relevantes para o adapter (Sprint 1–3)

Todos padrão Ethereum, sem extensão proprietária conhecida para os casos de uso do MVP:

| Operação no `chainport.Adapter` | Método(s) JSON-RPC |
|---|---|
| `CreateWallet` | Não é uma chamada RPC — é geração local de par de chaves secp256k1 e derivação do endereço (ver decisão abaixo) |
| `GetBalance` | `eth_getBalance` |
| `SendTransaction` | `eth_sendRawTransaction` (a transação é assinada localmente antes de ser enviada) |
| `EstimateGas` | `eth_estimateGas` |
| (suporte) | `eth_chainId` — usado para validar que o client está conectado à rede esperada e para assinatura EIP-155 |

## Decisão: `CreateWallet` gera chave localmente (self-custody), não usa Circle Wallets

A documentação da Arc aponta para o produto **Circle Wallets** (wallet-as-a-service, custódia gerenciada pela Circle) como caminho recomendado para lidar com carteiras. Optamos por **não** usar esse produto no adapter: gerar o par de chaves localmente (secp256k1 padrão EVM) mantém o adapter portável para qualquer chain EVM futura sem depender de uma API proprietária da Circle, e é consistente com o princípio *Blockchain Agnostic* do charter — `chainport.Adapter.CreateWallet` deve significar a mesma coisa (uma chave que o Aureon controla) em toda rede suportada. Ver decisão formal em [tradeoffs.md](../tradeoffs.md#tr-007-cliente-rpc-da-arc-via-go-ethereum-e-gera%C3%A7%C3%A3o-local-de-chave-em-vez-de-circle-wallets).

## Diferenças conhecidas em relação a um EVM genérico

A documentação oficial (`docs.arc.io`) menciona que existem diferenças em relação ao EVM padrão, mas não as detalha na página consultada nesta pesquisa. **Ação de acompanhamento**: revisitar `https://docs.arc.io/arc/references/gas-and-fees` antes de Sprint 3 (EstimateGas/SendTransaction), já que o gas nativo em USDC pode implicar em regras de precificação de gas diferentes do EVM padrão (ex: conversão de unidade, mínimos).

## Fontes

- [RPC endpoints — Arc Docs](https://docs.arc.io/arc/references/rpc-endpoints)
- [Arc — visão geral (docs.arc.io/llms.txt)](https://docs.arc.io/llms.txt)
- [Announcing Support for Arc's Public Testnet — Blockdaemon](https://www.blockdaemon.com/blog/announcing-support-for-arcs-public-testnet)
- [Circle Launches Arc Public Testnet — Circle](https://www.circle.com/pressroom/circle-launches-arc-public-testnet)
