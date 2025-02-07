package db

import (
	"context"

	"github.com/moth13/finance_tracker/util"
	decimal "github.com/shopspring/decimal"
)

type addMoneyTxParams struct {
	Amount      decimal.Decimal
	FinalAmount decimal.Decimal
	AccountID   int64
	MonthID     int64
	YearID      int64
}

func addMoneyTx(ctx context.Context,
	q *Queries,
	arg addMoneyTxParams) (balance util.Balance, err error) {

	// Update final balance for account, month and years
	argAccount := AddAccountBalanceParams{
		ID:          arg.AccountID,
		Amount:      arg.Amount,
		FinalAmount: arg.FinalAmount,
	}
	account, err := q.AddAccountBalance(ctx, argAccount)
	if err != nil {
		return
	}

	argMonth := AddMonthBalanceParams{
		ID:          arg.MonthID,
		Amount:      arg.Amount,
		FinalAmount: arg.FinalAmount,
	}
	month, err := q.AddMonthBalance(ctx, argMonth)
	if err != nil {
		return
	}

	argYear := AddYearBalanceParams{
		ID:          arg.YearID,
		Amount:      arg.Amount,
		FinalAmount: arg.FinalAmount,
	}
	year, err := q.AddYearBalance(ctx, argYear)
	if err != nil {
		return
	}

	balance.AccountBalance = account.Balance
	balance.AccountFinalBalance = account.FinalBalance
	balance.MonthBalance = month.Balance
	balance.MonthFinalBalance = month.FinalBalance
	balance.YearBalance = year.Balance
	balance.YearFinalBalance = year.FinalBalance

	return balance, nil
}
