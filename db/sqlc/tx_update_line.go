package db

import (
	"context"
	"fmt"
	"time"

	"github.com/moth13/finance_tracker/util"
	decimal "github.com/shopspring/decimal"
)

// UpdateLineTxParams contains all infos to update a new line
type UpdateLineTxParams struct {
	ID          int64               `json:"id"`
	Title       *string             `json:"title"`
	Owner       *string             `json:"owner"`
	AccountID   *int64              `json:"account_id"`
	MonthID     *int64              `json:"month_id"`
	YearID      *int64              `json:"year_id"`
	CategoryID  *int64              `json:"category_id"`
	Amount      decimal.NullDecimal `json:"amount"`
	Checked     *bool               `json:"checked"`
	Description *string             `json:"description"`
	DueDate     *time.Time          `json:"due_date"`
}

// UpdateLineTxResult contains all infos about the result of line creation
type UpdateLineTxResult struct {
	Line    Line         `json:"line"`
	Balance util.Balance `json:"balance"`
}

func (store *SQLStore) UpdateLineTx(ctx context.Context, arg UpdateLineTxParams) (UpdateLineTxResult, error) {
	var result UpdateLineTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Check if line exists
		line, err := q.GetLine(ctx, arg.ID)
		if err != nil {
			return err
		}

		argLine := UpdateLineParams{
			ID:          line.ID,
			Title:       line.Title,
			Description: line.Description,
			Checked:     line.Checked,
			Amount:      line.Amount,
			AccountID:   line.AccountID,
			MonthID:     line.MonthID,
			YearID:      line.YearID,
			CategoryID:  line.CategoryID,
			DueDate:     line.DueDate,
		}

		// Overload when needs it
		if arg.Title != nil {
			argLine.Title = *arg.Title
		}

		if arg.Description != nil {
			argLine.Description = *arg.Description
		}

		if arg.Checked != nil {
			argLine.Checked = *arg.Checked
		}

		if arg.Amount.Valid {
			argLine.Amount = arg.Amount.Decimal
		}

		if arg.CategoryID != nil {
			argLine.CategoryID = *arg.CategoryID
		}

		if arg.AccountID != nil {
			argLine.AccountID = *arg.AccountID
		}

		if arg.MonthID != nil {
			argLine.MonthID = *arg.MonthID
		}

		if arg.YearID != nil {
			argLine.YearID = *arg.YearID
		}

		if arg.DueDate != nil {
			argLine.DueDate = *arg.DueDate
		}

		// Revert previous balance for all components
		argRevert := addMoneyTxParams{
			Amount:      decimal.Zero,
			FinalAmount: line.Amount.Neg(),
			AccountID:   line.AccountID,
			MonthID:     line.MonthID,
			YearID:      line.YearID,
		}

		if line.Checked {
			argRevert.Amount = line.Amount.Neg()
		}
		result.Balance, err = addMoneyTx(ctx, q, argRevert)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Apply new balance
		argUpdate := addMoneyTxParams{
			Amount:      decimal.Zero,
			FinalAmount: argLine.Amount,
			AccountID:   argLine.AccountID,
			MonthID:     argLine.MonthID,
			YearID:      argLine.YearID,
		}

		if line.Checked {
			argUpdate.Amount = argLine.Amount
		}
		result.Balance, err = addMoneyTx(ctx, q, argUpdate)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Update the line
		result.Line, err = q.UpdateLine(ctx, argLine)
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
