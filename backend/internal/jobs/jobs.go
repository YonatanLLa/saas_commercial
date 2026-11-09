// Package jobs registra los trabajos en segundo plano de River (BR-117).
// Cada trabajo nuevo (procesar webhook, enviar mensaje, ejecutar automatización,
// recordatorio de cita…) se define con un Args y un Worker y se registra en Register.
package jobs

import (
	"context"
	"log/slog"

	"github.com/riverqueue/river"
)

// PingArgs es un trabajo de ejemplo para comprobar que la cola funciona.
type PingArgs struct {
	Message string `json:"message"`
}

func (PingArgs) Kind() string { return "ping" }

type PingWorker struct {
	river.WorkerDefaults[PingArgs]
	Logger *slog.Logger
}

func (w *PingWorker) Work(_ context.Context, job *river.Job[PingArgs]) error {
	w.Logger.Info("ping recibido", "job_id", job.ID, "message", job.Args.Message)
	return nil
}

// Register agrega todos los workers. Lo usan el binario worker y las pruebas.
func Register(logger *slog.Logger) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &PingWorker{Logger: logger})
	return workers
}
