package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(*testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("before tx:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	//check the result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// TODO: Check acc balance
		fmt.Println("tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("after tx:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTwoWay(t *testing.T) {
	store := NewStore(*testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 10
	amount := int64(10)

	errs := make(chan error)
	fmt.Println("before tx:", account1.Balance, account2.Balance)

	for i := 0; i < n; i++ {
		arg := TransferTxParam{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        amount,
		}
		if i%2 == 0 {
			arg.FromAccountID = account2.ID
			arg.ToAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), arg)
			errs <- err
		}()
	}

	//check the result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("after tx:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
