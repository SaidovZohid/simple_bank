package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/SaidovZohid/simple_bank/pkg/utils"
	"gitlab.com/SaidovZohid/simple_bank/storage/repo"
)

func createRandomTransfer(t *testing.T, account1, account2 *repo.Account) *repo.Transfer {
	arg := repo.Transfer{
		FromAccountId: account1.Id,
		ToAccountId:   account2.Id,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := dbManager.Transfer().Create(context.Background(), &arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountId, transfer.FromAccountId)
	require.Equal(t, arg.ToAccountId, transfer.ToAccountId)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.Id)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createAccount(t)
	account2 := createAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createAccount(t)
	account2 := createAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := dbManager.Transfer().Get(context.Background(), transfer1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.Id, transfer2.Id)
	require.Equal(t, transfer1.FromAccountId, transfer2.FromAccountId)
	require.Equal(t, transfer1.ToAccountId, transfer2.ToAccountId)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	account1 := createAccount(t)
	account2 := createAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := repo.TransferParams{
		FromAccountId: account1.Id,
		ToAccountId:   account1.Id,
		Limit:         5,
		Page:          1,
	}

	transfers, err := dbManager.Transfer().GetAll(context.Background(), &arg)
	require.NoError(t, err)
	require.Len(t, transfers.Transfers, 5)

	for _, transfer := range transfers.Transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountId == account1.Id || transfer.ToAccountId == account1.Id)
	}
}
