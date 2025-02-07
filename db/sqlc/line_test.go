package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomLine(t *testing.T, user User, month Month, year Year, account Account, category Category) Line {
	arg := CreateLineParams{
		Title:       util.RandomTitle(),
		Owner:       user.Username,
		AccountID:   account.ID,
		MonthID:     month.ID,
		YearID:      year.ID,
		CategoryID:  category.ID,
		Amount:      util.RandomMoney(),
		Checked:     util.RandomBool(),
		DueDate:     util.RandomFutureDate(),
		Description: util.RandomString(14),
	}

	line, err := testStore.CreateLine(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, line)

	require.NotZero(t, line.ID)
	require.Equal(t, line.Owner, arg.Owner)
	require.Equal(t, line.Title, arg.Title)
	require.Equal(t, line.AccountID, arg.AccountID)
	require.Equal(t, line.MonthID, arg.MonthID)
	require.True(t, line.Amount.Equal(arg.Amount))
	require.Equal(t, line.Checked, arg.Checked)
	require.Equal(t, line.Description, arg.Description)
	require.Equal(t, line.YearID, arg.YearID)
	require.Equal(t, line.CategoryID, arg.CategoryID)
	require.WithinDuration(t, line.DueDate, arg.DueDate, time.Second)

	return line
}

func TestCreateLine(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	createRandomLine(t, user, month, year, account, category)
}

func TestGetLine(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	line1 := createRandomLine(t, user, month, year, account, category)

	line2, err := testStore.GetLine(context.Background(), line1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, line2)

	require.NotZero(t, line2.ID)
	require.Equal(t, line1.Owner, line2.Owner)
	require.Equal(t, line1.Title, line2.Title)
	require.Equal(t, line1.AccountID, line2.AccountID)
	require.Equal(t, line1.MonthID, line2.MonthID)
	require.Equal(t, line1.YearID, line2.YearID)
	require.Equal(t, line1.CategoryID, line2.CategoryID)
	require.True(t, line1.Amount.Equal(line2.Amount))
	require.Equal(t, line1.Checked, line2.Checked)
	require.Equal(t, line1.Description, line2.Description)
	require.WithinDuration(t, line1.DueDate, line2.DueDate, time.Second)
}

func TestDeleteLine(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	line1 := createRandomLine(t, user, month, year, account, category)

	err := testStore.DeleteLine(context.Background(), line1.ID)
	require.NoError(t, err)

	line2, err := testStore.GetLine(context.Background(), line1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, line2)
}

func TestUpdateLine(t *testing.T) {
	user := createRandomUser(t)
	account1 := createRandomAccount(t, user)
	year1 := createRandomYear(t, user)
	month1 := createRandomMonth(t, user, year1)
	category1 := createRandomCategory(t, user)

	line1 := createRandomLine(t, user, month1, year1, account1, category1)

	account2 := createRandomAccount(t, user)
	year2 := createRandomYear(t, user)
	month2 := createRandomMonth(t, user, year2)
	category2 := createRandomCategory(t, user)

	arg := UpdateLineParams{
		ID:          line1.ID,
		Title:       util.RandomTitle(),
		Description: util.RandomString(14),
		AccountID:   account2.ID,
		MonthID:     month2.ID,
		YearID:      year2.ID,
		CategoryID:  category2.ID,
		Amount:      util.RandomMoney(),
		Checked:     util.RandomBool(),
		DueDate:     util.RandomFutureDate(),
	}

	line2, err := testStore.UpdateLine(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, line2)

	require.NotZero(t, line2.ID)
	require.Equal(t, line2.AccountID, arg.AccountID)
	require.Equal(t, line2.MonthID, arg.MonthID)
	require.Equal(t, line2.YearID, arg.YearID)
	require.Equal(t, line2.CategoryID, arg.CategoryID)
	require.Equal(t, line2.Title, arg.Title)
	require.Equal(t, line2.Checked, arg.Checked)
	require.Equal(t, line2.Description, arg.Description)
	require.WithinDuration(t, line2.DueDate, arg.DueDate, time.Second)
	require.True(t, line2.Amount.Equal(arg.Amount))
}

func TestListLines(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	year := createRandomYear(t, user)
	month := createRandomMonth(t, user, year)
	category := createRandomCategory(t, user)

	var lastLine Line
	for i := 0; i < 10; i++ {
		lastLine = createRandomLine(t, user, month, year, account, category)
	}

	arg := ListLinesParams{
		Owner:  lastLine.Owner,
		Limit:  5,
		Offset: 0,
	}

	lines, err := testStore.ListLines(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, lines)

	for _, line := range lines {
		require.NotEmpty(t, line)
		require.Equal(t, lastLine.Owner, line.Owner)
	}
}
