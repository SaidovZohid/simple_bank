package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
	"gitlab.com/SaidovZohid/simple_bank/pkg/utils"
	"gitlab.com/SaidovZohid/simple_bank/storage/repo"
)

func createAccount(t *testing.T) *repo.Account {
	model := &repo.Account{
		Owner:    faker.FirstName(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := dbManager.Account().Create(context.Background(), model)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, model.Owner, account.Owner)
	require.Equal(t, model.Balance, account.Balance)
	require.Equal(t, model.Currency, account.Currency)

	require.NotZero(t, account.Id)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	account := createAccount(t)
	err := dbManager.Account().Delete(context.Background(), account.Id)
	require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	account1 := createAccount(t)
	account2, err := dbManager.Account().Get(context.Background(), account1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Id, account2.Id)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	err = dbManager.Account().Delete(context.Background(), account1.Id)
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createAccount(t)

	arg := &repo.Account{
		Id:      account1.Id,
		Balance: utils.RandomMoney(),
	}

	account2, err := dbManager.Account().Update(context.Background(), arg.Id, arg.Balance)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Id, account2.Id)
	require.Equal(t, account1.Owner, account2.Owner)
	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	err = dbManager.Account().Delete(context.Background(), account1.Id)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createAccount(t)
	err := dbManager.Account().Delete(context.Background(), account1.Id)
	require.NoError(t, err)

	account2, err := dbManager.Account().Get(context.Background(), account1.Id)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestGetAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createAccount(t)
	}

	arg := &repo.AccountParams{
		Limit: 10,
		Page:  1,
	}

	accounts, err := dbManager.Account().GetAll(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Limit, int64(len(accounts.Accounts)))
	for _, account := range accounts.Accounts {
		require.NotEmpty(t, account)
		err = dbManager.Account().Delete(context.Background(), account.Id)
		require.NoError(t, err)
	}
}
