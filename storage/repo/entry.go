package repo

import (
	"context"
	"time"
)

type EntryStorageI interface {
	Create(ctx context.Context, c *Entry) (*Entry, error)
	Get(ctx context.Context, entryId int64) (*Entry, error)
	GetAll(ctx context.Context, params *EntryParams) (*Entries, error)
}

type Entry struct {
	Id        int64
	AccountId int64
	// can be positive and negative
	Amount    float64
	CreatedAt time.Time
}

type Entries struct {
	Entries []*Entry
	Count   int64
}

type EntryParams struct {
	Limit     int64
	Page      int64
	AccountId int64
}
