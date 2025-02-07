package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomYear(t *testing.T, user User) Year {
	start, end := util.RandomYearDate()
	arg := CreateYearParams{
		Title:       util.RandomTitle(),
		Owner:       user.Username,
		Description: util.RandomString(14),
		StartDate:   start,
		EndDate:     end,
	}

	year, err := testStore.CreateYear(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, year)

	require.NotZero(t, year.ID)
	require.Equal(t, year.Title, arg.Title)
	require.Equal(t, year.Description, arg.Description)
	require.WithinDuration(t, year.StartDate, arg.StartDate, time.Second)
	require.WithinDuration(t, year.EndDate, arg.EndDate, time.Second)
	require.True(t, year.Balance.IsZero())
	require.True(t, year.FinalBalance.IsZero())

	return year
}

func TestCreateYear(t *testing.T) {
	user := createRandomUser(t)
	createRandomYear(t, user)
}

func TestGetYear(t *testing.T) {
	user := createRandomUser(t)

	year1 := createRandomYear(t, user)

	year2, err := testStore.GetYear(context.Background(), year1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, year2)

	require.NotZero(t, year2.ID)
	require.Equal(t, year2.Title, year1.Title)
	require.Equal(t, year2.Description, year1.Description)
	require.WithinDuration(t, year2.StartDate, year1.StartDate, time.Second)
	require.WithinDuration(t, year2.EndDate, year1.EndDate, time.Second)
	require.True(t, year2.Balance.IsZero())
	require.True(t, year2.FinalBalance.IsZero())
}

func TestUpdateYear(t *testing.T) {
	start, end := util.RandomYearDate()
	user := createRandomUser(t)

	year1 := createRandomYear(t, user)

	arg := UpdateYearParams{
		ID:          year1.ID,
		Title:       util.RandomTitle(),
		Description: util.RandomString(14),
		StartDate:   start,
		EndDate:     end,
	}

	year2, err := testStore.UpdateYear(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, year2)

	require.NotZero(t, year2.ID)
	require.Equal(t, year2.Title, arg.Title)
	require.Equal(t, year2.Description, arg.Description)
	require.WithinDuration(t, year2.StartDate, arg.StartDate, time.Second)
	require.WithinDuration(t, year2.EndDate, arg.EndDate, time.Second)
	require.True(t, year2.Balance.Equal(year1.Balance))
	require.True(t, year2.FinalBalance.Equal(year1.FinalBalance))
}

func TestAddYearBalance(t *testing.T) {
	user := createRandomUser(t)

	year1 := createRandomYear(t, user)

	arg := AddYearBalanceParams{
		ID:          year1.ID,
		Amount:      util.RandomMoney(),
		FinalAmount: util.RandomMoney(),
	}

	year2, err := testStore.AddYearBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, year2)

	require.Equal(t, year1.ID, year2.ID)
	require.Equal(t, year1.Owner, year2.Owner)
	require.Equal(t, year1.Title, year2.Title)
	require.True(t, year2.Balance.Equal(year1.Balance.Add(arg.Amount)))
	require.WithinDuration(t, year1.StartDate, year2.StartDate, time.Second)
	require.WithinDuration(t, year1.EndDate, year2.EndDate, time.Second)
	require.True(t, year2.FinalBalance.Equal(year1.FinalBalance.Add(arg.FinalAmount)))
}

func TestDeleteYear(t *testing.T) {
	user := createRandomUser(t)

	year1 := createRandomYear(t, user)

	err := testStore.DeleteYear(context.Background(), year1.ID)
	require.NoError(t, err)

	year2, err := testStore.GetYear(context.Background(), year1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, year2)
}

func TestListYears(t *testing.T) {
	user := createRandomUser(t)

	var lastYear Year
	for i := 0; i < 10; i++ {
		lastYear = createRandomYear(t, user)
	}

	arg := ListYearsParams{
		Owner:  lastYear.Owner,
		Limit:  5,
		Offset: 0,
	}

	years, err := testStore.ListYears(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, years)

	for _, year := range years {
		require.NotEmpty(t, year)
		require.Equal(t, lastYear.Owner, year.Owner)
	}
}
