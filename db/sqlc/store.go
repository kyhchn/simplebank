package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db sql.DB) *Store {
	return &Store{
		db:      &db,
		Queries: New(&db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := New(tx)
	err = fn(query)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx errorr: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// function to exexute transafer transaction containing create transfer, entry for 2 account, and edit each account
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry1")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry2")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}
		//get each account and run update query
		if arg.FromAccountID > arg.ToAccountID {
			transferMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)
		} else {
			transferMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})
	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = q.UpdateBalance(ctx, UpdateBalanceParams{
		Amount: -amount,
		ID:     fromAccountID,
	})
	if err != nil {
		return
	}

	toAccount, err = q.UpdateBalance(ctx, UpdateBalanceParams{
		Amount: amount,
		ID:     toAccountID,
	})
	if err != nil {
		return
	}
	return
}
