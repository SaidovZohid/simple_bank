package repo

import (
	"context"
	"time"
)

type AccountStorageI interface {
	Create(ctx context.Context, a *Account) (*Account, error)
	Get(ctx context.Context, accountId int64) (*Account, error)
	Update(ctx context.Context, accountId int64, balance float64) (*Account, error)
	Delete(ctx context.Context, accountId int64) error
	GetAll(ctx context.Context, params *AccountParams) (*Accounts, error)
}

type Account struct {
	Id        int64
	Owner     string
	Balance   float64
	Currency  string
	CreatedAt time.Time
}

type Accounts struct {
	Accounts []*Account
	Count    int64
}

type AccountParams struct {
	Limit  int64
	Page   int64
	Owner  string
	SortBy string
}
