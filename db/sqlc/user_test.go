package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: "secret",
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
		Currency:       util.RandomCurrency(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.Currency, arg.Currency)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreateAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestCreateSameUsername(t *testing.T) {
	user1 := createRandomUser(t)
	arg := CreateUserParams{
		Username:       user1.Username,
		HashedPassword: "secret",
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
		Currency:       util.RandomCurrency(),
	}

	_, err := testStore.CreateUser(context.Background(), arg)
	require.Error(t, err)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	db, err := testStore.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, db)

	require.Equal(t, user.Username, db.Username)
	require.Equal(t, user.HashedPassword, db.HashedPassword)
	require.Equal(t, user.FullName, db.FullName)
	require.Equal(t, user.Email, db.Email)
	require.Equal(t, user.Currency, db.Currency)
	require.Equal(t, user.PasswordChangedAt, db.PasswordChangedAt)
	require.Equal(t, user.CreateAt, db.CreateAt)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testStore.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	account2, err := testStore.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}
