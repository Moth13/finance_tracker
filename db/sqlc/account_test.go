package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, user User) Account {
	arg := CreateAccountParams{
		Owner:       user.Username,
		Title:       util.RandomTitle(),
		InitBalance: util.RandomMoney(),
		Description: util.RandomString(10),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.NotZero(t, account.ID)
	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Title, arg.Title)
	require.True(t, account.InitBalance.Equal(arg.InitBalance))
	require.True(t, account.Balance.IsZero())
	require.True(t, account.FinalBalance.IsZero())

	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)

	createRandomAccount(t, user)
}

func TestGetAccount(t *testing.T) {
	user := createRandomUser(t)

	account := createRandomAccount(t, user)

	db, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, db)

	require.NotZero(t, account.ID)
	require.Equal(t, account.Owner, db.Owner)
	require.Equal(t, account.Title, db.Title)
	require.True(t, account.InitBalance.Equal(db.InitBalance))
	require.True(t, account.Balance.Equal(db.Balance))
	require.True(t, account.FinalBalance.Equal(db.FinalBalance))
}

func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)

	account1 := createRandomAccount(t, user)

	arg := UpdateAccountParams{
		ID:          account1.ID,
		InitBalance: util.RandomMoney(),
		Title:       util.RandomTitle(),
		Description: util.RandomString(14),
	}

	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account2.Title, arg.Title)
	require.Equal(t, account2.Description, arg.Description)
	require.Equal(t, account2.Title, arg.Title)
	require.True(t, arg.InitBalance.Equal(account2.InitBalance))
	require.True(t, account1.Balance.Equal(account2.Balance))
	require.True(t, account1.FinalBalance.Equal(account2.FinalBalance))
}

func TestAddAccountBalance(t *testing.T) {
	user := createRandomUser(t)

	account1 := createRandomAccount(t, user)

	arg := AddAccountBalanceParams{
		ID:          account1.ID,
		Amount:      util.RandomMoney(),
		FinalAmount: util.RandomMoney(),
	}

	account2, err := testStore.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Title, account2.Title)
	require.True(t, account1.InitBalance.Equal(account2.InitBalance))
	require.True(t, account2.Balance.Equal(account1.Balance.Add(arg.Amount)))
	require.True(t, account2.FinalBalance.Equal(account1.FinalBalance.Add(arg.FinalAmount)))
}

func TestDeleteAccount(t *testing.T) {
	user := createRandomUser(t)

	account1 := createRandomAccount(t, user)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestDeleteInvalidAccount(t *testing.T) {
	user := createRandomUser(t)

	account1 := createRandomAccount(t, user)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	user := createRandomUser(t)

	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t, user)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
