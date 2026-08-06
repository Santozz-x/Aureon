package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeielsantos/aureon/modules/chains/arc"
	arcrpc "github.com/jeielsantos/aureon/modules/chains/arc/rpc"
	chainport "github.com/jeielsantos/aureon/modules/contracts"
	"github.com/jeielsantos/aureon/modules/platform/internal/adapter/rest"
	"github.com/jeielsantos/aureon/modules/platform/internal/infra/config"
	"github.com/jeielsantos/aureon/modules/platform/internal/infra/keystore"
	"github.com/jeielsantos/aureon/modules/platform/internal/usecase"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	bootCtx, cancelBoot := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelBoot()

	arcClient, err := arcrpc.Dial(bootCtx, cfg.ArcRPCURL)
	if err != nil {
		logger.Error("failed to connect to ARC Network RPC", "url", cfg.ArcRPCURL, "error", err)
		os.Exit(1)
	}
	defer arcClient.Close()

	keys := keystore.NewMemory()

	adapters := map[chainport.Network]chainport.Adapter{
		chainport.Network("arc"): arc.NewAdapter(arcClient, keys),
	}

	walletService := usecase.NewWalletService(adapters)
	walletHandler := rest.NewWalletHandler(walletService)

	transactionService := usecase.NewTransactionService(adapters)
	transactionHandler := rest.NewTransactionHandler(transactionService)

	router := rest.NewRouter(walletHandler, transactionHandler)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		logger.Info("starting aureon gateway", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("gateway server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down aureon gateway")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
