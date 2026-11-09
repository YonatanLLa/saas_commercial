// Binario del worker: procesa los trabajos de River en segundo plano
// (webhooks, envíos de WhatsApp, automatizaciones, recordatorios, IA).
// Se despliega aparte de la API para escalar por separado (RNF-002).
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"saas/backend/internal/jobs"
	"saas/backend/internal/platform/config"
	"saas/backend/internal/platform/database"
	"saas/backend/internal/platform/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("el worker terminó con error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New(cfg.Env, cfg.LogLevel)
	slog.SetDefault(log)

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	client, err := river.NewClient(riverpgxv5.New(db), &river.Config{
		Logger: log,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 50},
		},
		Workers: jobs.Register(log),
	})
	if err != nil {
		return err
	}

	if err := client.Start(ctx); err != nil {
		return err
	}
	log.Info("worker iniciado")

	<-ctx.Done()
	log.Info("apagando worker…")

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return client.Stop(stopCtx)
}
