package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomRecLine(t *testing.T, user User, account Account, category Category) Recline {
	arg := CreateRecLineParams{
		Title:       util.RandomTitle(),
		Owner:       user.Username,
		AccountID:   account.ID,
		CategoryID:  category.ID,
		Amount:      util.RandomMoney(),
		DueDate:     util.RandomFutureDate(),
		Recurrency:  util.RandomRecurrency(),
		Description: util.RandomString(14),
	}

	recline, err := testStore.CreateRecLine(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, recline)

	require.NotZero(t, recline.ID)
	require.Equal(t, recline.Title, arg.Title)
	require.Equal(t, recline.AccountID, arg.AccountID)
	require.True(t, recline.Amount.Equal(arg.Amount))
	require.Equal(t, recline.Recurrency, arg.Recurrency)
	require.Equal(t, recline.Description, arg.Description)
	require.Equal(t, recline.CategoryID, arg.CategoryID)
	require.WithinDuration(t, recline.DueDate, arg.DueDate, time.Second)

	return recline
}

func TestCreateRecLine(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	category := createRandomCategory(t, user)

	createRandomRecLine(t, user, account, category)
}

func TestDeleteRecLine(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	category := createRandomCategory(t, user)

	line1 := createRandomRecLine(t, user, account, category)

	err := testStore.DeleteRecLine(context.Background(), line1.ID)
	require.NoError(t, err)

	line2, err := testStore.GetRecLine(context.Background(), line1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, line2)
}

func TestGetRecLine(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	category := createRandomCategory(t, user)

	line1 := createRandomRecLine(t, user, account, category)

	line2, err := testStore.GetRecLine(context.Background(), line1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, line2)

	require.NotZero(t, line2.ID)
	require.Equal(t, line1.Owner, line2.Owner)
	require.Equal(t, line1.Title, line2.Title)
	require.Equal(t, line1.AccountID, line2.AccountID)
	require.Equal(t, line1.CategoryID, line2.CategoryID)
	require.True(t, line1.Amount.Equal(line2.Amount))
	require.Equal(t, line1.Description, line2.Description)
	require.WithinDuration(t, line1.DueDate, line2.DueDate, time.Second)
}

func TestListRecLines(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
	category := createRandomCategory(t, user)

	var lastLine Recline
	for i := 0; i < 10; i++ {
		lastLine = createRandomRecLine(t, user, account, category)
	}

	arg := ListRecLinesParams{
		Owner:  lastLine.Owner,
		Limit:  5,
		Offset: 0,
	}

	lines, err := testStore.ListRecLines(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, lines)

	for _, line := range lines {
		require.NotEmpty(t, line)
		require.Equal(t, lastLine.Owner, line.Owner)
	}
}
