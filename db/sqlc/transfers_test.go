package db

import (
	"context"
	"testing"

	"github.com/kyhchn/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	transfer, err := testQuery.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)
	require.Equal(t, transfer.Amount, arg.Amount)
	return transfer
}
func TestCreateTransfer(t *testing.T) {
	account := createRandomTransfer(t)
	require.NotEmpty(t, account)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQuery.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
}

func TestListTransfers(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	for i := 0; i < 10; i++ {
		dumpAccount := createRandomAccount(t)
		arg := CreateTransferParams{
			FromAccountID: dumpAccount.ID,
			ToAccountID:   transfer1.ToAccountID,
			Amount:        util.RandomBalance(),
		}
		transfer, err := testQuery.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)

		require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
		require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
		require.Equal(t, transfer.Amount, arg.Amount)
	}
	arg := ListTransfersParams{
		ToAccountID: transfer1.ToAccountID,
		Limit:       5,
		Offset:      5,
	}
	transfers, err := testQuery.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
