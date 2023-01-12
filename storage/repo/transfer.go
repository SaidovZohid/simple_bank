package repo

import (
	"context"
	"time"
)

type TransferStorageI interface {
	Create(ctx context.Context, c *Transfer) (*Transfer, error)
	Get(ctx context.Context, transferId int64) (*Transfer, error)
	GetAll(ctx context.Context, params *TransferParams) (*Transfers, error)
}

type Transfer struct {
	Id            int64
	FromAccountId int64
	ToAccountId   int64
	// must be postive
	Amount    float64
	CreatedAt time.Time
}

type Transfers struct {
	Transfers []*Transfer
	Count     int64
}

type TransferParams struct {
	Limit         int64
	Page          int64
	FromAccountId int64
	ToAccountId   int64
	SortBy        string
}
