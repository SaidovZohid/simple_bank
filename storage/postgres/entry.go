package postgres

import (
	"context"
	"fmt"

	"gitlab.com/SaidovZohid/simple_bank/storage/repo"

	"github.com/jmoiron/sqlx"
)

type entryRepo struct {
	db *sqlx.DB
}

func NewEntryStorage(db *sqlx.DB) repo.EntryStorageI {
	return &entryRepo{
		db: db,
	}
}

func (d *entryRepo) Create(ctx context.Context, entry *repo.Entry) (*repo.Entry, error) {
	query := `
		INSERT INTO entries (
			account_id,
			amount
		) VALUES ($1, $2) 
		RETURNING id, created_at
	`
	err := d.db.QueryRow(
		query,
		entry.AccountId,
		entry.Amount,
	).Scan(
		&entry.Id,
		&entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (d *entryRepo) Get(ctx context.Context, entryId int64) (*repo.Entry, error) {
	var (
		entry repo.Entry
	)
	query := `
		SELECT
			id,
			account_id,
			amount,
			created_at
		FROM entries WHERE id = $1	
	`
	err := d.db.QueryRow(
		query,
		entryId,
	).Scan(
		&entry.Id,
		&entry.AccountId,
		&entry.Amount,
		&entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (d *entryRepo) GetAll(ctx context.Context, params *repo.EntryParams) (*repo.Entries, error) {
	var (
		entries repo.Entries
	)
	entries.Entries = make([]*repo.Entry, 0)

	offset := (params.Page - 1) * params.Limit
	filter := ""
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	if params.AccountId > 0 {
		filter = fmt.Sprintf(" WHERE account_id = %d ", params.AccountId)
	}
	query := `
		SELECT
			id,
			account_id,
			amount,
			created_at
		FROM entries	
	` + filter + " ORDER BY id " + limit
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e repo.Entry
		err := rows.Scan(
			&e.Id,
			&e.AccountId,
			&e.Amount,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		entries.Entries = append(entries.Entries, &e)
	}

	return &entries, nil
}
