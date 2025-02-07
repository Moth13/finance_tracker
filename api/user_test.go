package api

import (
	"testing"

	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
		Currency:       util.RandomCurrency(),
	}
	return
}
