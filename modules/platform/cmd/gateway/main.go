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

	"github.com/Santozz-x/Aureon/modules/chains/arc"
	arcrpc "github.com/Santozz-x/Aureon/modules/chains/arc/rpc"
	chainport "github.com/Santozz-x/Aureon/modules/contracts"
	"github.com/Santozz-x/Aureon/modules/platform/internal/adapter/rest"
	"github.com/Santozz-x/Aureon/modules/platform/internal/adapter/rest/middleware"
	"github.com/Santozz-x/Aureon/modules/platform/internal/infra/apikeystore"
	"github.com/Santozz-x/Aureon/modules/platform/internal/infra/config"
	"github.com/Santozz-x/Aureon/modules/platform/internal/infra/db"
	"github.com/Santozz-x/Aureon/modules/platform/internal/infra/keystore"
	"github.com/Santozz-x/Aureon/modules/platform/internal/usecase"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	bootCtx, cancelBoot := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelBoot()

	arcClient, err := arcrpc.Dial(bootCtx, cfg.ArcRPCURL)
	if err != nil {
		logger.Error("failed to connect to ARC Network RPC", "url", cfg.ArcRPCURL, "error", err)
		os.Exit(1)
	}
	defer arcClient.Close()

	sqlDB, err := db.Open(bootCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	if err := db.Migrate(sqlDB); err != nil {
		logger.Error("failed to apply database migrations", "error", err)
		os.Exit(1)
	}

	keys, err := keystore.NewPostgres(sqlDB, cfg.KeystoreEncryptionKey)
	if err != nil {
		logger.Error("failed to init key store", "error", err)
		os.Exit(1)
	}

	adapters := map[chainport.Network]chainport.Adapter{
		chainport.Network("arc"): arc.NewAdapter(arcClient, keys),
	}

	walletService := usecase.NewWalletService(adapters)
	walletHandler := rest.NewWalletHandler(walletService)

	transactionService := usecase.NewTransactionService(adapters)
	transactionHandler := rest.NewTransactionHandler(transactionService)

	apiKeyService := usecase.NewAPIKeyService(apikeystore.NewPostgres(sqlDB))
	apiKeyHandler := rest.NewAPIKeyHandler(apiKeyService)
	protect := middleware.RequireAPIKey(apiKeyService)

	router := rest.NewRouter(walletHandler, transactionHandler, apiKeyHandler, protect)
	handler := middleware.Logging(logger)(router)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
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
