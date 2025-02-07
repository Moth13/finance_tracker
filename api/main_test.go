package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
