package db

import (
	"context"

	"github.com/moth13/finance_tracker/util"
	decimal "github.com/shopspring/decimal"
)

// DeleteLineTxParams contains all infos to create a new line
type DeleteLineTxParams struct {
	ID int64 `json:"id"`
}

// AddLineTxResult contains all infos about the result of line creation
type DeleteLineTxResult struct {
	Balance util.Balance `json:"balance"`
}

func (store *SQLStore) DeleteLineTx(ctx context.Context, arg DeleteLineTxParams) (DeleteLineTxResult, error) {
	var result DeleteLineTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		line, err := q.GetLine(ctx, arg.ID)
		if err != nil {
			return err
		}

		argAdd := addMoneyTxParams{
			Amount: decimal.Zero,
			FinalAmount: line.Amount.Neg(),
			AccountID: line.AccountID,
			MonthID: line.MonthID,
			YearID: line.YearID,
		}

		if line.Checked {
			argAdd.Amount = line.Amount.Neg()
		}

		// Update balance for each parts, ie substract the amount of the line
		result.Balance, err = addMoneyTx(ctx, q, argAdd)
		if err != nil {
			return err
		}

		err = q.DeleteLine(ctx, arg.ID)
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
