package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/SaidovZohid/simple_bank/pkg/utils"
	"gitlab.com/SaidovZohid/simple_bank/storage/repo"
)

func createRandomEntry(t *testing.T, account *repo.Account) *repo.Entry {
	arg := repo.Entry{
		AccountId: account.Id,
		Amount:    utils.RandomMoney(),
	}

	entry, err := dbManager.Entry().Create(context.Background(), &arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountId, entry.AccountId)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.Id)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := dbManager.Entry().Get(context.Background(), entry1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.Id, entry2.Id)
	require.Equal(t, entry1.AccountId, entry2.AccountId)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := repo.EntryParams{
		Limit:     10,
		Page:      1,
		AccountId: account.Id,
	}

	entries, err := dbManager.Entry().GetAll(context.Background(), &arg)
	require.NoError(t, err)
	require.Len(t, entries.Entries, 10)

	for _, entry := range entries.Entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountId, entry.AccountId)
	}
}
