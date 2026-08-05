# Contribuindo

Aureon é open source (Apache 2.0) desde o início. Este guia cobre o essencial para contribuir com o monorepo.

## Antes de começar

1. Leia o [CHARTER.md](../CHARTER.md) para entender a visão e os princípios do projeto.
2. Veja o [implementation-plan.md](implementation-plan.md) para saber o que está planejado e em qual sprint/fase.
3. Para mudanças de arquitetura, registre a decisão em [tradeoffs.md](tradeoffs.md) (o que foi escolhido, alternativas, o que se sacrifica).

## Ambiente

- Go 1.26+
- `go work sync` após clonar, para sincronizar os módulos do `go.work`

## Fluxo de trabalho

```bash
make build   # compila todos os módulos
make vet     # go vet em todos os módulos
make test    # testes de todos os módulos
make fmt     # gofmt -s -w em todo o repositório
```

Todo PR deve passar em `make build`, `make vet` e `make test` antes de ser aberto.

## Convenções

- **Módulos independentes por blockchain**: uma nova rede é sempre um novo módulo em `modules/chains/{network}`, implementando `chainport.Adapter` (definido em `modules/contracts`). Nunca adicione lógica específica de rede em `modules/platform`.
- **Camadas do `platform`**: `internal/domain` e `internal/usecase` nunca importam um módulo de chain específico — apenas a porta em `modules/contracts`. Wiring concreto só acontece em `cmd/*/main.go`. Ver [ARCHITECTURE.md](ARCHITECTURE.md).
- **Sem comentários óbvios**: comente apenas o "porquê" não-óbvio (uma decisão, uma restrição escondida), nunca o "o quê" — nomes de função/variável já devem explicar isso.
- **Commits**: mensagens no imperativo, descrevendo a motivação da mudança, não só o que mudou.

## Reportando problemas

Abra uma issue descrevendo o comportamento esperado vs. observado. Para vulnerabilidades de segurança, não abra uma issue pública — veja `SECURITY.md` (a ser criado antes do primeiro release público).
