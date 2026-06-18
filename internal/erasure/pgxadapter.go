// pgxadapter.go — pgxpool-backed implementation of the erasure.DB + Tx
// abstractions. Splitting the adapter from the executor keeps the executor's
// unit tests free of any pgx import and lets cmd/erasure-worker wire the real
// pool only at the boundary.
package erasure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxDB wraps a pgxpool.Pool to satisfy DB. Each BeginTx returns a fresh Tx
// that has already issued `SET TRANSACTION PRIORITY HIGH` — guaranteeing the
// erasure mutation wins contention against routine workloads (CockroachDB
// honours the priority for the entire transaction).
type PgxDB struct{ pool *pgxpool.Pool }

// NewPgxDB constructs the adapter. The pool is shared (the caller owns close).
func NewPgxDB(pool *pgxpool.Pool) *PgxDB { return &PgxDB{pool: pool} }

func (a *PgxDB) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := a.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, fmt.Errorf("pgx.BeginTx: %w", err)
	}
	if _, err := tx.Exec(ctx, `SET TRANSACTION PRIORITY HIGH`); err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("pgx.SetPriority: %w", err)
	}
	return &pgxTx{tx: tx}, nil
}

// pgxTx adapts pgx.Tx to the local Tx interface — only the surface the
// Executor uses (Exec / Commit / Rollback).
type pgxTx struct{ tx pgx.Tx }

func (t *pgxTx) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	tag, err := t.tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (t *pgxTx) Commit(ctx context.Context) error   { return t.tx.Commit(ctx) }
func (t *pgxTx) Rollback(ctx context.Context) error { return t.tx.Rollback(ctx) }
