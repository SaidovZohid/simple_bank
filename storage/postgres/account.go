package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/SaidovZohid/simple_bank/storage/repo"

	"github.com/jmoiron/sqlx"
)

type accountRepo struct {
	db *sqlx.DB
}

func NewAccountStorage(db *sqlx.DB) repo.AccountStorageI {
	return &accountRepo{
		db: db,
	}
}

func (d *accountRepo) Create(ctx context.Context, account *repo.Account) (*repo.Account, error) {
	query := `
		INSERT INTO accounts (
			owner,
			balance,
			currency
		) VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	err := d.db.QueryRow(
		query,
		account.Owner,
		account.Balance,
		account.Currency,
	).Scan(
		&account.Id,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (d *accountRepo) Get(ctx context.Context, accountId int64) (*repo.Account, error) {
	var (
		account repo.Account
	)
	query := `
		SELECT
			id,
			owner,
			balance,
			currency,
			created_at
		FROM accounts WHERE id = $1	
	`
	err := d.db.QueryRow(
		query,
		accountId,
	).Scan(
		&account.Id,
		&account.Owner,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (d *accountRepo) Update(ctx context.Context, accountId int64, balance float64) (*repo.Account, error) {
	var (
		res repo.Account
	)
	query := `
		UPDATE accounts SET 
			balance = $1
		WHERE id = $2 
		RETURNING id, owner, balance, currency, created_at
	`
	err := d.db.QueryRow(
		query,
		balance,
		accountId,
	).Scan(
		&res.Id,
		&res.Owner,
		&res.Balance,
		&res.Currency,
		&res.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *accountRepo) Delete(ctx context.Context, accountId int64) error {
	query := `
		DELETE FROM accounts WHERE id = $1
	`
	result, err := d.db.Exec(
		query,
		accountId,
	)
	if err != nil {
		return err
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *accountRepo) GetAll(ctx context.Context, params *repo.AccountParams) (*repo.Accounts, error) {
	var (
		accounts repo.Accounts
	)
	accounts.Accounts = make([]*repo.Account, 0)

	offset := (params.Page - 1) * params.Limit
	filter := ""
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	if params.Owner != "" {
		filter = fmt.Sprintf("WHERE owner ILIKE '%s'", "%"+params.Owner+"%")
	}
	orderBy := " ORDER BY created_at DESC "
	if params.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s ", params.SortBy)
	}
	query := `
		SELECT
			id,
			owner,
			balance,
			currency,
			created_at
		FROM accounts	
	` + filter + orderBy + limit
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a repo.Account
		err := rows.Scan(
			&a.Id,
			&a.Owner,
			&a.Balance,
			&a.Currency,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		accounts.Accounts = append(accounts.Accounts, &a)
	}

	return &accounts, nil
}
