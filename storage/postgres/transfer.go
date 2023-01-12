package postgres

import (
	"context"
	"fmt"

	"gitlab.com/SaidovZohid/simple_bank/storage/repo"

	"github.com/jmoiron/sqlx"
)

type transferRepo struct {
	db *sqlx.DB
}

func NewTransferStorage(db *sqlx.DB) repo.TransferStorageI {
	return &transferRepo{
		db: db,
	}
}

func (d *transferRepo) Create(ctx context.Context, transfer *repo.Transfer) (*repo.Transfer, error) {
	tx, err := d.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO transfers (
			from_account_id,
			to_account_id,
			amount
		) VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err = tx.QueryRow(
		query,
		transfer.FromAccountId,
		transfer.ToAccountId,
		transfer.Amount,
	).Scan(
		&transfer.Id,
		&transfer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (d *transferRepo) Get(ctx context.Context, tranferId int64) (*repo.Transfer, error) {
	var (
		tranfer repo.Transfer
	)
	query := `
		SELECT
			id,
			from_account_id,
			to_account_id,
			amount,
			created_at
		FROM transfers 
		WHERE id = $1
	`

	err := d.db.QueryRow(
		query,
		tranferId,
	).Scan(
		&tranfer.Id,
		&tranfer.FromAccountId,
		&tranfer.ToAccountId,
		&tranfer.Amount,
		&tranfer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tranfer, nil
}

func (d *transferRepo) GetAll(ctx context.Context, params *repo.TransferParams) (*repo.Transfers, error) {
	var (
		transfers repo.Transfers
	)
	transfers.Transfers = make([]*repo.Transfer, 0)

	offset := (params.Page - 1) * params.Limit
	filter := " WHERE true "
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	if params.FromAccountId > 0 || params.ToAccountId > 0 {
		filter += fmt.Sprintf("AND from_account_id = %d OR to_account_id = %d", params.FromAccountId, params.ToAccountId)
	}
	orderBy := " ORDER BY created_at DESC "
	if params.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s ", params.SortBy)
	}
	query := `
		SELECT
			id,
			from_account_id,
			to_account_id,
			amount,
			created_at
		FROM transfers	
	` + filter + orderBy + limit
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t repo.Transfer
		err := rows.Scan(
			&t.Id,
			&t.FromAccountId,
			&t.ToAccountId,
			&t.Amount,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		transfers.Transfers = append(transfers.Transfers, &t)
	}

	return &transfers, nil
}
