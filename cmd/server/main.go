package main

import (
	"log"

	"github.com/moth13/finance_tracker/api"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}

	conn, err := db.CreateDBConnection(config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't start server:", err)
	}
}
