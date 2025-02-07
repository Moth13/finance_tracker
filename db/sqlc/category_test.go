package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T, user User) Category {
	arg := CreateCategoryParams{
		Title: util.RandomTitle(),
		Owner: user.Username,
	}

	category, err := testStore.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.NotZero(t, category.ID)
	require.Equal(t, category.Title, arg.Title)
	require.Equal(t, category.Owner, arg.Owner)

	return category
}

func TestCreateCategory(t *testing.T) {
	user := createRandomUser(t)
	createRandomCategory(t, user)
}

func TestGetCategory(t *testing.T) {
	user := createRandomUser(t)

	category1 := createRandomCategory(t, user)

	category2, err := testStore.GetCategory(context.Background(), category1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.NotZero(t, category2.ID)
	require.Equal(t, category2.Title, category1.Title)
	require.Equal(t, category2.Owner, category1.Owner)
}

func TestDeleteCategory(t *testing.T) {
	user := createRandomUser(t)

	category1 := createRandomCategory(t, user)

	err := testStore.DeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	category2, err := testStore.GetCategory(context.Background(), category1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, category2)
}

func TestListCategories(t *testing.T) {
	user := createRandomUser(t)

	var lastCategory Category
	for i := 0; i < 10; i++ {
		lastCategory = createRandomCategory(t, user)
	}

	arg := ListCategoriesParams{
		Owner:  lastCategory.Owner,
		Limit:  5,
		Offset: 0,
	}

	categories, err := testStore.ListCategories(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	for _, category := range categories {
		require.NotEmpty(t, category)
		require.Equal(t, lastCategory.Owner, category.Owner)
	}
}
