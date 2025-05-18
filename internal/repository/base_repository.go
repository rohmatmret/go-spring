package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-spring.com/internal/observability"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db        *sql.DB
	txManager *TransactionManager
	tracer    *observability.Tracer
	tableName string
	idColumn  string
	columns   []string
	columnMap map[string]string
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB, tableName string, idColumn string, columns []string) *BaseRepository {
	columnMap := make(map[string]string)
	for _, col := range columns {
		columnMap[col] = col
	}

	return &BaseRepository{
		db:        db,
		txManager: NewTransactionManager(db),
		tracer:    observability.NewTracer("repository"),
		tableName: tableName,
		idColumn:  idColumn,
		columns:   columns,
		columnMap: columnMap,
	}
}

// FindByID finds an entity by its ID
func (r *BaseRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	start := time.Now()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		joinColumns(r.columns), r.tableName, r.idColumn)

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find by ID: %w", err)
	}

	observability.ServiceMethodDuration.WithLabelValues("Repository", "FindByID").Observe(time.Since(start).Seconds())
	return user, nil
}

// Create creates a new entity
func (r *BaseRepository) Create(ctx context.Context, entity interface{}) error {
	start := time.Now()
	tx, err := r.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Build insert query
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		r.tableName,
		joinColumns(r.columns),
		buildPlaceholders(len(r.columns)),
		r.idColumn)

	// Execute query
	err = tx.QueryRow(query, entity).Scan(&entity)
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	observability.ServiceMethodDuration.WithLabelValues("Repository", "Create").Observe(time.Since(start).Seconds())
	return nil
}

// Update updates an existing entity
func (r *BaseRepository) Update(ctx context.Context, id interface{}, entity interface{}) error {
	start := time.Now()
	tx, err := r.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Build update query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d",
		r.tableName,
		buildUpdateSet(r.columns),
		r.idColumn,
		len(r.columns)+1)

	// Execute query
	_, err = tx.Exec(query, append([]interface{}{entity}, id))
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	observability.ServiceMethodDuration.WithLabelValues("Repository", "Update").Observe(time.Since(start).Seconds())
	return nil
}

// Delete deletes an entity by its ID
func (r *BaseRepository) Delete(ctx context.Context, id interface{}) error {
	start := time.Now()
	tx, err := r.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", r.tableName, r.idColumn)
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	observability.ServiceMethodDuration.WithLabelValues("Repository", "Delete").Observe(time.Since(start).Seconds())
	return nil
}

// Helper functions
func joinColumns(columns []string) string {
	return fmt.Sprintf("%s", columns)
}

func buildPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf("%s", placeholders)
}

func buildUpdateSet(columns []string) string {
	sets := make([]string, len(columns))
	for i, col := range columns {
		sets[i] = fmt.Sprintf("%s = $%d", col, i+1)
	}
	return fmt.Sprintf("%s", sets)
}
