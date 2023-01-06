package db

import (
	"context"
	"testing"

	"github.com/kyhchn/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQuery.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotNil(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}
	account2, err := testQuery.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQuery.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, arr := testQuery.GetAccount(context.Background(), account.ID)
	require.Error(t, arr)
	require.Empty(t, account2)

}
