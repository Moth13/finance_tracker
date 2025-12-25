package util

import decimal "github.com/shopspring/decimal"

type Balance struct {
	MonthBalance        decimal.Decimal `json:"month_balance"`
	MonthFinalBalance   decimal.Decimal `json:"month_final_balance"`
	YearBalance         decimal.Decimal `json:"year_balance"`
	YearFinalBalance    decimal.Decimal `json:"year_final_balance"`
	AccountBalance      decimal.Decimal `json:"account_balance"`
	AccountFinalBalance decimal.Decimal `json:"account_final_balance"`
}
