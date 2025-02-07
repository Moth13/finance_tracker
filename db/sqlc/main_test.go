package db

import (
	"log"
	"os"
	"testing"

	"github.com/moth13/finance_tracker/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}

	connPool, err := CreateDBConnection(config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
