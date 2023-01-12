package storage

import (
	"gitlab.com/SaidovZohid/simple_bank/storage/postgres"

	"github.com/jmoiron/sqlx"
	"gitlab.com/SaidovZohid/simple_bank/storage/repo"
)

type StorageI interface {
	Account() repo.AccountStorageI
	Entry() repo.EntryStorageI
	Transfer() repo.TransferStorageI
}

type StoragePg struct {
	accountRepo  repo.AccountStorageI
	entryRepo    repo.EntryStorageI
	transferRepo repo.TransferStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &StoragePg{
		accountRepo:  postgres.NewAccountStorage(db),
		entryRepo:    postgres.NewEntryStorage(db),
		transferRepo: postgres.NewTransferStorage(db),
	}
}

func (s *StoragePg) Account() repo.AccountStorageI {
	return s.accountRepo
}

func (s *StoragePg) Entry() repo.EntryStorageI {
	return s.entryRepo
}

func (s *StoragePg) Transfer() repo.TransferStorageI {
	return s.transferRepo
}
