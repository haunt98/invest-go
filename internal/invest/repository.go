package invest

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	preparedCreateInvestment = "createInvestment"
	preparedDeleteInvestment = "deleteInvestment"

	stmtInitInvestments = `
	CREATE TABLE investments
	(
		id     TEXT PRIMARY KEY,
		amount INTEGER NOT NULL,
		date   TEXT    NOT NULL,
		source TEXT    NOT NULL
	)
`

	stmtGetInvestments = `
	SELECT id, amount, date, source
	FROM investments
`

	stmtCreateInvestment = `
	INSERT INTO investments (id, amount, date, source)
	VALUES (?, ?, ?, ?)
`

	stmtDeleteInvestment = `
	DELETE
	FROM investments
	WHERE id = ?
`
)

type Repository interface {
	GetInvestments(ctx context.Context) ([]Investment, error)
	CreateInvestment(ctx context.Context, investment Investment) error
	DeleteInvestment(ctx context.Context, id string) error
}

type repo struct {
	db *sql.DB

	// Prepared statements
	// https://go.dev/doc/database/prepared-statements
	preparedStmts map[string]*sql.Stmt
}

func NewRepository(ctx context.Context, db *sql.DB, shouldInitDatabase bool) (Repository, error) {
	if shouldInitDatabase {
		if _, err := db.ExecContext(ctx, stmtInitInvestments); err != nil {
			return nil, fmt.Errorf("database failed to exec: %w", err)
		}
	}

	preparedStmts := make(map[string]*sql.Stmt)

	var err error
	preparedStmts[preparedCreateInvestment], err = db.PrepareContext(ctx, stmtCreateInvestment)
	if err != nil {
		return nil, fmt.Errorf("database failed to prepare context: %w", err)
	}

	preparedStmts[preparedDeleteInvestment], err = db.PrepareContext(ctx, stmtDeleteInvestment)
	if err != nil {
		return nil, fmt.Errorf("database failed to prepare context: %w", err)
	}

	return &repo{
		db:            db,
		preparedStmts: preparedStmts,
	}, nil
}

func (r *repo) GetInvestments(ctx context.Context) ([]Investment, error) {
	investments := make([]Investment, 0, 16)

	rows, err := r.db.QueryContext(ctx, stmtGetInvestments)
	if err != nil {
		return nil, fmt.Errorf("database failed to query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		investment := Investment{}

		if err := rows.Scan(
			&investment.ID,
			&investment.Amount,
			&investment.Date,
			&investment.Source,
		); err != nil {
			return nil, fmt.Errorf("database failed to scan rows: %w", err)
		}

		investments = append(investments, investment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database failed to scan rows: %w", err)
	}

	return investments, nil
}

func (r *repo) CreateInvestment(ctx context.Context, investment Investment) error {
	if _, err := r.preparedStmts[preparedCreateInvestment].ExecContext(
		ctx,
		investment.ID,
		investment.Amount,
		investment.Date,
		investment.Source,
	); err != nil {
		return fmt.Errorf("database failed to exec: %w", err)
	}

	return nil
}

func (r *repo) DeleteInvestment(ctx context.Context, id string) error {
	if _, err := r.preparedStmts[preparedDeleteInvestment].ExecContext(
		ctx,
		id,
	); err != nil {
		return fmt.Errorf("database failed to exec: %w", err)
	}

	return nil
}
