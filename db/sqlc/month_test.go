package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomMonth(t *testing.T, user User, year Year) Month {
	start, end := util.RandomMonthDate()
	arg := CreateMonthParams{
		Title:       util.RandomTitle(),
		Owner:       user.Username,
		Description: util.RandomString(14),
		YearID:      year.ID,
		StartDate:   start,
		EndDate:     end,
	}

	month, err := testStore.CreateMonth(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, year)

	require.NotZero(t, month.ID)
	require.Equal(t, month.Title, arg.Title)
	require.Equal(t, month.Description, arg.Description)
	require.Equal(t, month.YearID, arg.YearID)
	require.WithinDuration(t, month.StartDate, arg.StartDate, time.Second)
	require.WithinDuration(t, month.EndDate, arg.EndDate, time.Second)
	require.True(t, month.Balance.IsZero())
	require.True(t, month.FinalBalance.IsZero())

	return month
}

func TestCreateMonth(t *testing.T) {
	user := createRandomUser(t)
	year := createRandomYear(t, user)
	createRandomMonth(t, user, year)
}

func TestAddMonthBalance(t *testing.T) {
	user := createRandomUser(t)
	year := createRandomYear(t, user)

	month1 := createRandomMonth(t, user, year)

	arg := AddMonthBalanceParams{
		ID:          month1.ID,
		Amount:      util.RandomMoney(),
		FinalAmount: util.RandomMoney(),
	}

	month2, err := testStore.AddMonthBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, month2)

	require.Equal(t, month1.ID, month2.ID)
	require.Equal(t, month1.Owner, month2.Owner)
	require.Equal(t, month1.Title, month2.Title)
	require.True(t, month2.Balance.Equal(month1.Balance.Add(arg.Amount)))
	require.WithinDuration(t, month1.StartDate, month2.StartDate, time.Second)
	require.WithinDuration(t, month1.EndDate, month2.EndDate, time.Second)
	require.True(t, month2.FinalBalance.Equal(month1.FinalBalance.Add(arg.FinalAmount)))
}

func TestDeleteMonth(t *testing.T) {
	user := createRandomUser(t)
	year := createRandomYear(t, user)

	month1 := createRandomMonth(t, user, year)

	err := testStore.DeleteMonth(context.Background(), month1.ID)
	require.NoError(t, err)

	month2, err := testStore.GetMonth(context.Background(), month1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, month2)
}

func TestUpdateMonth(t *testing.T) {
	start, end := util.RandomYearDate()
	user := createRandomUser(t)
	year1 := createRandomYear(t, user)
	year2 := createRandomYear(t, user)

	month1 := createRandomMonth(t, user, year1)

	arg := UpdateMonthParams{
		ID:          month1.ID,
		Title:       util.RandomTitle(),
		Description: util.RandomString(14),
		StartDate:   start,
		EndDate:     end,
		YearID:      year2.ID,
	}

	month2, err := testStore.UpdateMonth(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, month2)

	require.NotZero(t, month2.ID)
	require.Equal(t, month2.YearID, arg.YearID)
	require.Equal(t, month2.Title, arg.Title)
	require.Equal(t, month2.Description, arg.Description)
	require.WithinDuration(t, month2.StartDate, arg.StartDate, time.Second)
	require.WithinDuration(t, month2.EndDate, arg.EndDate, time.Second)
	require.True(t, month2.Balance.Equal(year1.Balance))
	require.True(t, month2.FinalBalance.Equal(year1.FinalBalance))
}

func TestListMonths(t *testing.T) {
	user := createRandomUser(t)
	year := createRandomYear(t, user)

	var lastMonth Month
	for i := 0; i < 10; i++ {
		lastMonth = createRandomMonth(t, user, year)
	}

	arg := ListMonthsParams{
		Owner:  lastMonth.Owner,
		Limit:  5,
		Offset: 0,
	}

	months, err := testStore.ListMonths(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, months)

	for _, month := range months {
		require.NotEmpty(t, month)
		require.Equal(t, lastMonth.Owner, month.Owner)
	}
}
