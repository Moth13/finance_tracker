package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	decimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestAddLineTx(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	// run n concurrent add line
	n := 5

	errs := make(chan error)
	results := make(chan AddLineTxResult)
	amounts := make(chan decimal.Decimal)
	for i := 0; i < n; i++ {
		go func() {
			tamount := util.RandomMoney()
			ctx := context.Background()
			result, err := testStore.AddLineTx(ctx, AddLineTxParams{
				Owner:       user.Username,
				Title:       util.RandomTitle(),
				Description: util.RandomString(14),
				Checked:     i%2 == 0,
				Amount:      tamount,
				AccountID:   account.ID,
				MonthID:     month.ID,
				YearID:      year.ID,
				CategoryID:  category.ID,
				DueDate:     time.Now(),
			})
			errs <- err
			results <- result
			amounts <- tamount
		}()
	}

	// Check results
	line_final_balance := decimal.Zero
	line_balance := decimal.Zero
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		amount := <-amounts
		require.NotEmpty(t, amount)

		// Check Line
		line := result.Line
		require.NotEmpty(t, line)
		require.NotZero(t, line.ID)
		require.Equal(t, line.AccountID, account.ID)
		require.Equal(t, line.MonthID, month.ID)
		require.Equal(t, line.YearID, year.ID)
		require.Equal(t, line.CategoryID, category.ID)
		require.True(t, line.Amount.Equal(amount))

		line_final_balance = line_final_balance.Add(line.Amount)
		if line.Checked {
			line_balance = line_balance.Add(line.Amount)
		}
	}

	// Check Account Balance
	updatedAccount, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	// fmt.Println("line_balance >>", line_balance)
	// fmt.Println("line_final_balance >>", line_final_balance)
	// fmt.Println("updatedAccount.Balance >>", updatedAccount.Balance)
	// fmt.Println("updatedAccount.FinalBalance >>", updatedAccount.FinalBalance)
	require.True(t, updatedAccount.Balance.Equal(account.Balance.Add(line_balance)))
	require.True(t, updatedAccount.FinalBalance.Equal(account.FinalBalance.Add(line_final_balance)))

	// Check Month Balance
	updateMonth, err := testStore.GetMonth(context.Background(), month.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateMonth)
	require.True(t, updateMonth.Balance.Equal(month.Balance.Add(line_balance)))
	require.True(t, updateMonth.FinalBalance.Equal(month.FinalBalance.Add(line_final_balance)))

	// Check Year Balance
	updateYear, err := testStore.GetYear(context.Background(), year.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateYear)
	require.True(t, updateYear.Balance.Equal(year.Balance.Add(line_balance)))
	require.True(t, updateYear.FinalBalance.Equal(year.FinalBalance.Add(line_final_balance)))
}

func TestDeleteLineTx(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	// run n concurrent add line
	n := 5

	errs := make(chan error)
	results := make(chan AddLineTxResult)
	amounts := make(chan decimal.Decimal)
	for i := 0; i < n; i++ {
		go func() {
			tamount := util.RandomMoney()
			ctx := context.Background()
			result, err := testStore.AddLineTx(ctx, AddLineTxParams{
				Owner:       user.Username,
				Title:       util.RandomTitle(),
				Description: util.RandomString(14),
				Checked:     i%2 == 0,
				Amount:      tamount,
				AccountID:   account.ID,
				MonthID:     month.ID,
				YearID:      year.ID,
				CategoryID:  category.ID,
				DueDate:     time.Now(),
			})
			errs <- err
			results <- result
			amounts <- tamount
		}()
	}

	// Check results
	line_final_balance := decimal.Zero
	line_balance := decimal.Zero
	lines := [5]Line{}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		amount := <-amounts
		require.NotEmpty(t, amount)

		// Check Line
		line := result.Line
		require.NotEmpty(t, line)
		require.NotZero(t, line.ID)
		require.Equal(t, line.AccountID, account.ID)
		require.Equal(t, line.MonthID, month.ID)
		require.Equal(t, line.YearID, year.ID)
		require.Equal(t, line.CategoryID, category.ID)
		require.True(t, line.Amount.Equal(amount))

		line_final_balance = line_final_balance.Add(line.Amount)
		if line.Checked {
			line_balance = line_balance.Add(line.Amount)
		}
		lines[i] = line
	}

	del_errs := make(chan error)
	del_results := make(chan DeleteLineTxResult)
	del_line := make(chan Line)
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := testStore.DeleteLineTx(ctx, DeleteLineTxParams{
				ID: lines[i].ID,
			})
			del_errs <- err
			del_results <- result
			del_line <- lines[i]
		}()
	}

	for i := 0; i < n; i++ {
		err := <-del_errs
		require.NoError(t, err)

		result := <-del_results
		require.NotEmpty(t, result)

		line := <-del_line
		require.NotEmpty(t, line)

		// Check Line deletion
		deletedLine, err := testStore.GetLine(context.Background(), line.ID)
		require.Error(t, err)
		require.EqualError(t, err, pgx.ErrNoRows.Error())
		require.Empty(t, deletedLine)

		line_final_balance = line_final_balance.Sub(line.Amount)
		if line.Checked {
			line_balance = line_balance.Sub(line.Amount)
		}
	}

	// Check Account Balance
	updatedAccount, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.True(t, updatedAccount.Balance.Equal(account.Balance.Add(line_balance)))
	require.True(t, updatedAccount.FinalBalance.Equal(account.FinalBalance.Add(line_final_balance)))

	// Check Month Balance
	updateMonth, err := testStore.GetMonth(context.Background(), month.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateMonth)
	require.True(t, updateMonth.Balance.Equal(month.Balance.Add(line_balance)))
	require.True(t, updateMonth.FinalBalance.Equal(month.FinalBalance.Add(line_final_balance)))

	// Check Year Balance
	updateYear, err := testStore.GetYear(context.Background(), year.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateYear)
	require.True(t, updateYear.Balance.Equal(year.Balance.Add(line_balance)))
	require.True(t, updateYear.FinalBalance.Equal(year.FinalBalance.Add(line_final_balance)))
}

func TestUpdateLineTx(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	// run n concurrent add line
	n := 5

	errs := make(chan error)
	results := make(chan AddLineTxResult)
	amounts := make(chan decimal.Decimal)
	for i := 0; i < n; i++ {
		go func() {
			tamount := util.RandomMoney()
			ctx := context.Background()
			result, err := testStore.AddLineTx(ctx, AddLineTxParams{
				Owner:       user.Username,
				Title:       util.RandomTitle(),
				Description: util.RandomString(14),
				Checked:     i%2 == 0,
				Amount:      tamount,
				AccountID:   account.ID,
				MonthID:     month.ID,
				YearID:      year.ID,
				CategoryID:  category.ID,
				DueDate:     time.Now(),
			})
			errs <- err
			results <- result
			amounts <- tamount
		}()
	}

	// Check results
	lines := [5]Line{}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		amount := <-amounts
		require.NotEmpty(t, amount)

		// Check Line
		line := result.Line
		require.NotEmpty(t, line)
		require.NotZero(t, line.ID)
		require.Equal(t, line.AccountID, account.ID)
		require.Equal(t, line.MonthID, month.ID)
		require.Equal(t, line.YearID, year.ID)
		require.Equal(t, line.CategoryID, category.ID)
		require.True(t, line.Amount.Equal(amount))

		lines[i] = line
	}

	up_errs := make(chan error)
	up_results := make(chan UpdateLineTxResult)
	up_line := make(chan Line)
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := testStore.UpdateLineTx(ctx, UpdateLineTxParams{
				ID:     lines[i].ID,
				Amount: decimal.NullDecimal{Decimal: util.RandomMoney(), Valid: true},
			})
			up_errs <- err
			up_results <- result
			up_line <- lines[i]
		}()
	}

	line_final_balance := decimal.Zero
	line_balance := decimal.Zero
	for i := 0; i < n; i++ {
		err := <-up_errs
		require.NoError(t, err)

		result := <-up_results
		require.NotEmpty(t, result)

		line := <-up_line
		require.NotEmpty(t, line)

		// Check Line update
		updatedLine, err := testStore.GetLine(context.Background(), line.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedLine)

		line_final_balance = line_final_balance.Add(updatedLine.Amount)
		if line.Checked {
			line_balance = line_balance.Add(updatedLine.Amount)
		}
	}

	// Check Account Balance
	updatedAccount, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.True(t, updatedAccount.Balance.Equal(account.Balance.Add(line_balance)))
	require.True(t, updatedAccount.FinalBalance.Equal(account.FinalBalance.Add(line_final_balance)))

	// Check Month Balance
	updateMonth, err := testStore.GetMonth(context.Background(), month.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateMonth)
	require.True(t, updateMonth.Balance.Equal(month.Balance.Add(line_balance)))
	require.True(t, updateMonth.FinalBalance.Equal(month.FinalBalance.Add(line_final_balance)))

	// Check Year Balance
	updateYear, err := testStore.GetYear(context.Background(), year.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateYear)
	require.True(t, updateYear.Balance.Equal(year.Balance.Add(line_balance)))
	require.True(t, updateYear.FinalBalance.Equal(year.FinalBalance.Add(line_final_balance)))
}
