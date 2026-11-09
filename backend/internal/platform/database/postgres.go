// Package database gestiona el pool de PostgreSQL y el contexto de tenant para RLS.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("creando pool de PostgreSQL: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("conectando a PostgreSQL: %w", err)
	}
	return pool, nil
}

// WithTenantTx ejecuta fn dentro de una transacción con app.tenant_id fijado.
// Las políticas RLS usan ese valor, de modo que ninguna consulta dentro de fn
// puede leer ni escribir filas de otro tenant (BR-010, BR-011).
// El tenantID debe venir de la membresía autenticada, nunca del cliente (BR-007).
func WithTenantTx(ctx context.Context, pool *pgxpool.Pool, tenantID string, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op si ya se hizo commit

	// El tercer argumento "true" limita el valor a esta transacción.
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID); err != nil {
		return fmt.Errorf("fijando tenant: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
