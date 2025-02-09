// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"time"

	decimal "github.com/shopspring/decimal"
)

type Account struct {
	ID           int64           `json:"id"`
	Owner        string          `json:"owner"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	InitBalance  decimal.Decimal `json:"init_balance"`
	Balance      decimal.Decimal `json:"balance"`
	FinalBalance decimal.Decimal `json:"final_balance"`
}

type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Owner string `json:"owner"`
}

type Line struct {
	ID         int64  `json:"id"`
	Owner      string `json:"owner"`
	Title      string `json:"title"`
	AccountID  int64  `json:"account_id"`
	MonthID    int64  `json:"month_id"`
	YearID     int64  `json:"year_id"`
	CategoryID int64  `json:"category_id"`
	// can be negative or positive
	Amount      decimal.Decimal `json:"amount"`
	Checked     bool            `json:"checked"`
	Description string          `json:"description"`
	DueDate     time.Time       `json:"due_date"`
}

type Month struct {
	ID           int64           `json:"id"`
	Owner        string          `json:"owner"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	YearID       int64           `json:"year_id"`
	Balance      decimal.Decimal `json:"balance"`
	FinalBalance decimal.Decimal `json:"final_balance"`
	StartDate    time.Time       `json:"start_date"`
	EndDate      time.Time       `json:"end_date"`
}

type Recline struct {
	ID        int64  `json:"id"`
	Owner     string `json:"owner"`
	Title     string `json:"title"`
	AccountID int64  `json:"account_id"`
	// can be negative or position
	Amount      decimal.Decimal `json:"amount"`
	CategoryID  int64           `json:"category_id"`
	Description string          `json:"description"`
	Recurrency  string          `json:"recurrency"`
	DueDate     time.Time       `json:"due_date"`
}

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Currency          string    `json:"currency"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreateAt          time.Time `json:"create_at"`
}

type Year struct {
	ID           int64           `json:"id"`
	Owner        string          `json:"owner"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Balance      decimal.Decimal `json:"balance"`
	FinalBalance decimal.Decimal `json:"final_balance"`
	StartDate    time.Time       `json:"start_date"`
	EndDate      time.Time       `json:"end_date"`
}
