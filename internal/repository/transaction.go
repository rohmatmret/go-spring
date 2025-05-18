package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-spring.com/internal/observability"
)

// TransactionManager handles database transactions
type TransactionManager struct {
	db     *sql.DB
	tracer *observability.Tracer
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db:     db,
		tracer: observability.NewTracer("transaction"),
	}
}

// Transaction represents a database transaction
type Transaction struct {
	tx     *sql.Tx
	ctx    context.Context
	tracer *observability.Tracer
}

// Begin starts a new transaction
func (tm *TransactionManager) Begin(ctx context.Context) (*Transaction, error) {
	start := time.Now()
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	observability.ServiceMethodDuration.WithLabelValues("Transaction", "Begin").Observe(time.Since(start).Seconds())

	return &Transaction{
		tx:     tx,
		ctx:    ctx,
		tracer: tm.tracer,
	}, nil
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	start := time.Now()
	err := t.tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	observability.ServiceMethodDuration.WithLabelValues("Transaction", "Commit").Observe(time.Since(start).Seconds())
	return nil
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	start := time.Now()
	err := t.tx.Rollback()
	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	observability.ServiceMethodDuration.WithLabelValues("Transaction", "Rollback").Observe(time.Since(start).Seconds())
	return nil
}

// Exec executes a query within the transaction
func (t *Transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := t.tx.ExecContext(t.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	observability.ServiceMethodDuration.WithLabelValues("Transaction", "Exec").Observe(time.Since(start).Seconds())
	return result, nil
}

// Query executes a query and returns rows
func (t *Transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := t.tx.QueryContext(t.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	observability.ServiceMethodDuration.WithLabelValues("Transaction", "Query").Observe(time.Since(start).Seconds())
	return rows, nil
}

// QueryRow executes a query and returns a single row
func (t *Transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := t.tx.QueryRowContext(t.ctx, query, args...)
	observability.ServiceMethodDuration.WithLabelValues("Transaction", "QueryRow").Observe(time.Since(start).Seconds())
	return row
}
