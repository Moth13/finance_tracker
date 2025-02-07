package db

import (
	"context"
	"time"

	"github.com/moth13/finance_tracker/util"
	decimal "github.com/shopspring/decimal"
)

// AddLineTxParams contains all infos to create a new line
type AddLineTxParams struct {
	Title       string          `json:"title"`
	Owner       string          `json:"owner"`
	Amount      decimal.Decimal `json:"amount"`
	Checked     bool            `json:"checked"`
	Description string          `json:"description"`
	DueDate     time.Time       `json:"due_date"`
	AccountID   int64           `json:"account_id"`
	MonthID     int64           `json:"month_id"`
	YearID      int64           `json:"year_id"`
	CategoryID  int64           `json:"category_id"`
}

// AddLineTxResult contains all infos about the result of line creation
type AddLineTxResult struct {
	Line    Line         `json:"line"`
	Balance util.Balance `json:"balance"`
}

func (store *SQLStore) AddLineTx(ctx context.Context, arg AddLineTxParams) (AddLineTxResult, error) {
	var result AddLineTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		argLine := CreateLineParams{
			Title:       arg.Title,
			Owner:       arg.Owner,
			Description: arg.Description,
			Checked:     arg.Checked,
			Amount:      arg.Amount,
			AccountID:   arg.AccountID,
			MonthID:     arg.MonthID,
			YearID:      arg.YearID,
			CategoryID:  arg.CategoryID,
			DueDate:     arg.DueDate,
		}

		argAdd := addMoneyTxParams{
			Amount: decimal.Zero,
			FinalAmount: arg.Amount,
			AccountID: arg.AccountID,
			MonthID: arg.MonthID,
			YearID: arg.YearID,
		}

		if arg.Checked {
			argAdd.Amount = arg.Amount
		}

		// Update balance for each parts, ie add the amount of the line
		result.Balance, err = addMoneyTx(ctx, q, argAdd)
		if err != nil {
			return err
		}

		if _, err := q.GetCategory(ctx, argLine.CategoryID); err != nil {
			return err
		}

		result.Line, err = q.CreateLine(ctx, argLine)
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
